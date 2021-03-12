package gitweb

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	gwProto "github.com/eelf/gitweb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"sync"
	"time"
)

const sessionCookie = "s"

type gitwebSvc struct {
	sessions sync.Map
}

type sessionTyp struct {
	userId string
	remoteAddr string
}

var gwService = &gitwebSvc{}

func RegisterService(grpcServer *grpc.Server) {
	gwProto.RegisterGitwebServer(grpcServer, gwService)
}

func (s *gitwebSvc) newSession() ([]byte, error) {
	sess := make([]byte, 24)
	session := sessionTyp{}
	for i := 0; i < 20; i++ {
		n, err := rand.Read(sess)
		if err != nil || n != len(sess) {
			return nil, fmt.Errorf("%v %v", err, n)
		}
		_, hit := s.sessions.LoadOrStore(string(sess), session)
		if !hit {
			return sess, nil
		}
	}
	return nil, fmt.Errorf("rand ne rand")
}

func (s *gitwebSvc) setAuth(ctx context.Context, id string) error {
	sess, err := s.newSession()

	md, _ := metadata.FromIncomingContext(ctx)
	fw := md.Get("x-forwarded-for")
	session := sessionTyp{
		userId: id,
	}
	if len(fw) > 0 {
		session.remoteAddr = fw[0]
	}

	if err != nil {
		return err
	}
	c := &http.Cookie{
		Name:    sessionCookie,
		Value:   base64.StdEncoding.EncodeToString(sess),
		Path:    "/",
		Expires: time.Now().Add(3e7*time.Second),
	}
	s.sessions.Store(c.Value, session)
	return grpc.SendHeader(ctx, metadata.Pairs("Set-Cookie", c.String()))
}

func (s *gitwebSvc) auth(ctx context.Context) (sessionTyp, error) {
	md, _ := metadata.FromIncomingContext(ctx)

	for _, cookie := range (&http.Request{
		Header: map[string][]string{"Cookie": md.Get("cookie")},
	}).Cookies() {
		if cookie.Name != sessionCookie {
			continue
		}
		id, ok := s.sessions.Load(cookie.Value)
		if ok {
			return id.(sessionTyp), nil
		}
	}
	return sessionTyp{}, fmt.Errorf("unauth")
}

func (s *gitwebSvc) RepoList(ctx context.Context, req *gwProto.RepoListRequest) (*gwProto.RepoListResponse, error) {
	userId, err := s.auth(ctx)
	if err != nil {
		return nil, err
	}
	log.Println(userId)

	var repos []*gwProto.RepoListResponseRepo
	for i := 0; i < 10; i++ {
		repos = append(repos, &gwProto.RepoListResponseRepo{
			Name: fmt.Sprint(i),
		})
	}
	return &gwProto.RepoListResponse{
		Repos: repos,
	}, nil
}

func (s *gitwebSvc) Auth(ctx context.Context, req *gwProto.AuthRequest) (resp *gwProto.AuthResponse, err error) {
	g := req.GetGoogleAccessToken()

	url := "https://www.googleapis.com/oauth2/v1/userinfo?access_token=" + g;
	var res *http.Response
	res, err = http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()
	googleId := struct {
		Id string `json:"id"`
	}{}
	err = json.NewDecoder(res.Body).Decode(&googleId)
	if err != nil {
		return
	}

	err = s.setAuth(ctx, googleId.Id)
	if err != nil {
		return
	}

	resp = &gwProto.AuthResponse{
		Text: googleId.Id,
	}

	return
}
