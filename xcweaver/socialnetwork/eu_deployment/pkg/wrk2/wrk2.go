package wrk2

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"eu_deployment/pkg/model"
	"eu_deployment/pkg/services"

	"github.com/XCWeaver/xcweaver"
)

type server struct {
	xcweaver.Implements[xcweaver.Main]
	xcweaver.WithConfig[serverOptions]
	homeTimelineService xcweaver.Ref[services.HomeTimelineService]
	userTimelineService xcweaver.Ref[services.UserTimelineService]
	textService         xcweaver.Ref[services.TextService]
	mediaService        xcweaver.Ref[services.MediaService]
	uniqueIdService     xcweaver.Ref[services.UniqueIdService]
	userService         xcweaver.Ref[services.UserService]
	socialGraphService  xcweaver.Ref[services.SocialGraphService]
	lis                 xcweaver.Listener `xcweaver:"wrk2"`
}

type serverOptions struct {
	Region string `toml:"region"`
}

func Serve(ctx context.Context, s *server) error {
	mux := http.NewServeMux()

	// declare api endpoints
	mux.Handle("/wrk2-api/user/register", instrument("user/register", s.registerHandler, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/user/follow", instrument("user/follow", s.followHandler, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/user/unfollow", instrument("user/unfollow", s.unfollowHandler, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/user/login", instrument("user/login", s.loginHandler, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/post/compose", instrument("post/compose", s.composePostHandler, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/post/update", instrument("post/update", s.updatePostHandler, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/home-timeline/read", instrument("home-timeline/read", s.readHomeTimelineHandler, http.MethodGet, http.MethodPost))
	mux.Handle("/wrk2-api/user-timeline/read", instrument("user-timeline/read", s.readUserTimelineHandler, http.MethodGet, http.MethodPost))

	var handler http.Handler = mux
	s.Logger(ctx).Info("wrk2-api available", "addr", s.lis, "region", s.Config().Region)
	return http.Serve(s.lis, handler)
}

func instrument(label string, fn func(http.ResponseWriter, *http.Request), methods ...string) http.Handler {
	allowed := map[string]struct{}{}
	for _, method := range methods {
		allowed[method] = struct{}{}
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		if _, ok := allowed[r.Method]; len(allowed) > 0 && !ok {
			msg := fmt.Sprintf("method %q not allowed", r.Method)
			http.Error(w, msg, http.StatusMethodNotAllowed)
			return
		}
		fn(w, r)
	}
	return xcweaver.InstrumentHandlerFunc(label, handler)
}

func genReqID() int64 {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
}

type registerParams struct {
	reqID     int64
	firstName string
	lastName  string
	username  string
	password  string
	userID    string
}

func validateRegisterParams(w http.ResponseWriter, r *http.Request) *registerParams {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "error: "+err.Error(), http.StatusBadRequest)
		return nil
	}

	params := registerParams{
		reqID:  genReqID(),
		userID: "",
	}
	// get params
	params.username = r.Form.Get("username")
	params.firstName = r.Form.Get("first_name")
	params.lastName = r.Form.Get("last_name")
	params.password = r.Form.Get("password")
	params.userID = r.Form.Get("user_id")

	// validate mandatory fields
	if params.username == "" {
		http.Error(w, "must provide a valid username", http.StatusBadRequest)
		return nil
	}
	if params.firstName == "" {
		http.Error(w, "must provide a valid first_name", http.StatusBadRequest)
		return nil
	}
	if params.lastName == "" {
		http.Error(w, "must provide a valid last_name", http.StatusBadRequest)
		return nil
	}

	return &params
}

func (s *server) registerHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering wkr2-api/user/register")

	params := validateRegisterParams(w, r)
	if params == nil {
		return
	}
	var err error
	if params.userID == "" {
		logger.Debug("calling userService.RegisterUser()", "reqID", params.reqID, "firstName", params.firstName, "lastName", params.lastName, "username", params.username, "password", params.password)
		err = s.userService.Get().RegisterUser(ctx, params.reqID, params.firstName, params.lastName, params.username, params.password)
	} else {
		logger.Debug("calling userService.RegisterUserWithId()", "reqID", params.reqID, "firstName", params.firstName, "lastName", params.lastName, "username", params.username, "password", params.password, "userID", params.userID)
		err = s.userService.Get().RegisterUserWithId(ctx, params.reqID, params.firstName, params.lastName, params.username, params.password, params.userID)
	}
	if err != nil {
		logger.Error("error registering user", "msg", err.Error())
		http.Error(w, "error registering user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Debug("success! registered user", "username", params.username, "userID", params.userID)
	response := fmt.Sprintf("success! registered user %s (id=%s)\n", params.username, params.userID)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(response))
}

type followParams struct {
	reqID        int64
	userID       string
	followeeID   string
	username     string
	followeeName string
}

func validateFollowParams(w http.ResponseWriter, r *http.Request) *followParams {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "error: "+err.Error(), http.StatusBadRequest)
		return nil
	}

	params := followParams{
		reqID:      genReqID(),
		userID:     "",
		followeeID: "",
	}
	// get params
	params.userID = r.Form.Get("user_id")
	params.followeeID = r.Form.Get("followee_id")
	params.username = r.Form.Get("user_name")
	params.followeeName = r.Form.Get("followee_name")

	return &params
}

func (s *server) followHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering wkr2-api/user/follow")

	var err error
	params := validateFollowParams(w, r)
	if params == nil {
		return
	}

	if params.userID != "" && params.followeeID != "" {
		logger.Debug("calling socialGraphService.Follow()", "reqID", params.reqID, "userID", params.userID, "followeeID", params.followeeID)
		err = s.socialGraphService.Get().Follow(ctx, params.reqID, params.userID, params.followeeID)
	} else if params.username != "" && params.followeeName != "" {
		logger.Debug("calling socialGraphService.FollowWithUsername()", "reqID", params.reqID, "username", params.username, "followeeName", params.followeeName)
		err = s.socialGraphService.Get().FollowWithUsername(ctx, params.reqID, params.username, params.followeeName)
	} else {
		logger.Error("error following user: invalid arguments", "userID", params.userID, "followeeID", params.followeeID, "username", params.username, "followeeName", params.followeeName)
		http.Error(w, "error following user: "+"invalid arguments", http.StatusBadRequest)
		return
	}

	if err != nil {
		logger.Error("error following user", "msg", err.Error())
		http.Error(w, "error following user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Debug("success! followed user", "follower username", params.username, "followed userID", params.userID, "followeeName", params.followeeName, "followeeID", params.followeeID)
	response := fmt.Sprintf("success! user %s (id=%d) followed user %s (id=%d)\n", params.username, params.userID, params.followeeName, params.followeeID)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(response))
}

func (s *server) unfollowHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering wkr2-api/user/unfollow")

	var err error
	params := validateFollowParams(w, r)
	if params == nil {
		return
	}
	if params.userID != "" && params.followeeID != "" {
		err = s.socialGraphService.Get().Unfollow(ctx, params.reqID, params.userID, params.followeeID)
	} else if params.username != "" && params.followeeName != "" {
		err = s.socialGraphService.Get().UnfollowWithUsername(ctx, params.reqID, params.username, params.followeeName)
	} else {
		logger.Error("error unfollowing user: invalid arguments", "userID", params.userID, "followeeID", params.followeeID, "username", params.username, "followeeName", params.followeeName)
		http.Error(w, "error unfollowing user: "+"invalid arguments", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "error unfollowing user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Debug("success! unfollowed user", "follower username", params.username, "followed userID", params.userID, "followeeName", params.followeeName, "followeeID", params.followeeID)
	response := fmt.Sprintf("success! user %s (id=%d) unfollowed user %s (id=%d)\n", params.username, params.userID, params.followeeName, params.followeeID)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(response))
}

type LoginParams struct {
	reqID    int64
	username string
	password string
}

func validateLoginParams(w http.ResponseWriter, r *http.Request) *LoginParams {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "error: "+err.Error(), http.StatusBadRequest)
		return nil
	}

	params := LoginParams{
		reqID: genReqID(),
	}
	// get params
	params.username = r.Form.Get("username")
	params.password = r.Form.Get("password")

	// validate mandatory fields
	if params.username == "" {
		http.Error(w, "must provide a valid username", http.StatusBadRequest)
		return nil
	}
	if params.password == "" {
		http.Error(w, "must provide a valid password", http.StatusBadRequest)
		return nil
	}
	return &params
}

func (s *server) loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering wkr2-api/user/login")

	params := validateLoginParams(w, r)
	if params == nil {
		return
	}

	token, err := s.userService.Get().Login(ctx, params.reqID, params.username, params.password)
	if err != nil {
		http.Error(w, "error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("success! user %s logged in with token %s\n", params.username, token)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(response))
}

type ComposePostParams struct {
	text       string
	userID     string
	username   string
	reqID      int64
	mediaTypes []string
	mediaIDs   []int64
	postType   model.PostType
}

type ComposePostReply struct {
	Text   string
	PostId string
}

func validateComposePostParams(logger *slog.Logger, r *http.Request) (*ComposePostParams, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	params := ComposePostParams{
		reqID:    genReqID(),
		userID:   "",
		postType: model.PostType(-1),
	}
	// get params
	params.username = r.Form.Get("username")
	params.text = r.Form.Get("text")
	params.userID = r.Form.Get("user_id")
	postTypeStr := r.Form.Get("post_type")
	mediaTypesStr := r.Form.Get("media_types")
	mediaIDsStr := r.Form.Get("media_ids")

	// validate types
	if postTypeStr != "" {
		postType, err := strconv.ParseInt(postTypeStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid post_type. Available types: 0-POST, 1-REPOST, 2-REPLY, 3-DM")
		}
		params.postType = model.PostType(postType)
	}
	if mediaTypesStr != "" && mediaTypesStr != "[]" {
		mediaTypesStr = strings.TrimPrefix(mediaTypesStr, "[")
		mediaTypesStr = strings.TrimSuffix(mediaTypesStr, "]")
		params.mediaTypes = strings.Split(mediaTypesStr, ",")
	}
	if mediaIDsStr != "" && mediaIDsStr != "[]" {
		mediaIDsStr = strings.TrimPrefix(mediaIDsStr, "[")
		mediaIDsStr = strings.TrimSuffix(mediaIDsStr, "]")
		mediaIDsStrSlice := strings.Split(mediaIDsStr, ",")
		for _, idStr := range mediaIDsStrSlice {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				return nil, fmt.Errorf(fmt.Sprintf("error parsing media ids: %s", err.Error()))
			}
			params.mediaIDs = append(params.mediaIDs, id)
		}
	}

	// validate mandatory fields
	if params.username == "" {
		return nil, fmt.Errorf("must provide a valid username")
	}
	if params.text == "" {
		return nil, fmt.Errorf("must provide a valid text")
	}
	if params.userID == "" {
		return nil, fmt.Errorf("must provide a user_id")
	}
	if params.postType < 0 || params.postType > 3 {
		return nil, fmt.Errorf("invalid post_type. Available types: 0-POST, 1-REPOST, 2-REPLY, 3-DM")
	}

	return &params, nil
}

func (s *server) composePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering wkr2-api/post/compose")

	params, err := validateComposePostParams(logger, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Debug("valid parameters", "params", params)

	var wg sync.WaitGroup
	wg.Add(4)
	var errs [4]error
	var id string
	go func() {
		defer wg.Done()
		logger.Debug("calling text service")
		errs[0] = s.textService.Get().UploadText(ctx, params.reqID, params.text, true)
		logger.Debug("upload text done!")
	}()
	go func() {
		defer wg.Done()
		logger.Debug("calling media service")
		errs[1] = s.mediaService.Get().UploadMedia(ctx, params.reqID, params.mediaTypes, params.mediaIDs, true)
		logger.Debug("upload media done!")
	}()
	go func() {
		defer wg.Done()
		logger.Debug("calling upload id service")
		id, errs[2] = s.uniqueIdService.Get().UploadUniqueId(ctx, params.reqID, params.postType, "")
		logger.Debug("upload unique id done!")
	}()
	go func() {
		defer wg.Done()
		logger.Debug("calling user service")
		errs[3] = s.userService.Get().UploadCreatorWithUserId(ctx, params.reqID, params.userID, params.username, true)
		logger.Debug("upload creator with user id done!")
	}()
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			logger.Debug("error composing post", "msg", err.Error())
			http.Error(w, "error composing post: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	logger.Debug("success! composed post", "username", params.username, "userID", params.userID, "text", params.text)
	response := fmt.Sprintf("success! user %s (id=%d) composed post: %s\n", params.username, params.userID, params.text)
	w.Header().Set("Content-Type", "application/json")
	replyJson, err := json.Marshal(ComposePostReply{response, id})
	if err != nil {
		logger.Error(err.Error())
		return
	}
	w.Write(replyJson)
}

type UpdatePostParams struct {
	postId     string
	text       string
	userID     string
	username   string
	reqID      int64
	mediaTypes []string
	mediaIDs   []int64
	postType   model.PostType
}

func validateUpdatePostParams(logger *slog.Logger, r *http.Request) (*UpdatePostParams, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	params := UpdatePostParams{
		reqID:    genReqID(),
		userID:   "",
		postType: model.PostType(-1),
	}
	// get params
	params.postId = r.Form.Get("post_id")
	params.username = r.Form.Get("username")
	params.text = r.Form.Get("text")
	params.userID = r.Form.Get("user_id")
	postTypeStr := r.Form.Get("post_type")
	mediaTypesStr := r.Form.Get("media_types")
	mediaIDsStr := r.Form.Get("media_ids")

	// validate types
	if postTypeStr != "" {
		postType, err := strconv.ParseInt(postTypeStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid post_type. Available types: 0-POST, 1-REPOST, 2-REPLY, 3-DM")
		}
		params.postType = model.PostType(postType)
	}
	if mediaTypesStr != "" && mediaTypesStr != "[]" {
		mediaTypesStr = strings.TrimPrefix(mediaTypesStr, "[")
		mediaTypesStr = strings.TrimSuffix(mediaTypesStr, "]")
		params.mediaTypes = strings.Split(mediaTypesStr, ",")
	}
	if mediaIDsStr != "" && mediaIDsStr != "[]" {
		mediaIDsStr = strings.TrimPrefix(mediaIDsStr, "[")
		mediaIDsStr = strings.TrimSuffix(mediaIDsStr, "]")
		mediaIDsStrSlice := strings.Split(mediaIDsStr, ",")
		for _, idStr := range mediaIDsStrSlice {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				return nil, fmt.Errorf(fmt.Sprintf("error parsing media ids: %s", err.Error()))
			}
			params.mediaIDs = append(params.mediaIDs, id)
		}
	}

	// validate mandatory fields
	if params.username == "" {
		return nil, fmt.Errorf("must provide a valid username")
	}
	if params.text == "" {
		return nil, fmt.Errorf("must provide a valid text")
	}
	if params.userID == "" {
		return nil, fmt.Errorf("must provide a user_id")
	}
	if params.postId == "" {
		return nil, fmt.Errorf("must provide a post_id")
	}
	if params.postType < 0 || params.postType > 3 {
		return nil, fmt.Errorf("invalid post_type. Available types: 0-POST, 1-REPOST, 2-REPLY, 3-DM")
	}

	return &params, nil
}

func (s *server) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering wkr2-api/post/update")

	params, err := validateUpdatePostParams(logger, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Debug("valid parameters", "params", params)

	var wg sync.WaitGroup
	wg.Add(4)
	var errs [4]error
	go func() {
		defer wg.Done()
		logger.Debug("calling text service")
		errs[0] = s.textService.Get().UploadText(ctx, params.reqID, params.text, false)
		logger.Debug("upload text done!")
	}()
	go func() {
		defer wg.Done()
		logger.Debug("calling media service")
		errs[1] = s.mediaService.Get().UploadMedia(ctx, params.reqID, params.mediaTypes, params.mediaIDs, false)
		logger.Debug("upload media done!")
	}()
	go func() {
		defer wg.Done()
		logger.Debug("calling upload id service")
		_, errs[2] = s.uniqueIdService.Get().UploadUniqueId(ctx, params.reqID, params.postType, params.postId)
		logger.Debug("upload unique id done!")
	}()
	go func() {
		defer wg.Done()
		logger.Debug("calling user service")
		errs[3] = s.userService.Get().UploadCreatorWithUserId(ctx, params.reqID, params.userID, params.username, false)
		logger.Debug("upload creator with user id done!")
	}()
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			logger.Debug("error composing post", "msg", err.Error())
			http.Error(w, "error composing post: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	logger.Debug("success! post updated", "username", params.username, "userID", params.userID, "text", params.text)
	response := fmt.Sprintf("success! user %s (id=%d) updated post: %s\n", params.username, params.userID, params.text)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(response))
}

type ReadTimelineParams struct {
	reqID  int64
	userID string
	start  int64
	stop   int64
}

func validateReadTimelineParams(w http.ResponseWriter, r *http.Request) *ReadTimelineParams {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "error: "+err.Error(), http.StatusBadRequest)
		return nil
	}

	params := ReadTimelineParams{
		reqID:  genReqID(),
		userID: "",
		start:  0,
		stop:   10,
	}
	// get params
	params.userID = r.Form.Get("user_id")

	// validate types
	if params.userID == "" {
		http.Error(w, "must provide a user_id", http.StatusBadRequest)
		return nil
	}

	return &params
}

func (s *server) readHomeTimelineHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering wkr2-api/home-timeline/read")

	params := validateReadTimelineParams(w, r)
	if params == nil {
		return
	}
	posts, err := s.homeTimelineService.Get().ReadHomeTimeline(ctx, params.reqID, params.userID, params.start, params.stop)
	if err != nil {
		http.Error(w, "error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// response, err := json.Marshal(posts)
	// iterate to add new line in between posts
	for i, post := range posts {
		postJSON, err := json.Marshal(post)
		if err != nil {
			http.Error(w, "error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(postJSON)
		if i < len(posts)-1 {
			w.Write([]byte("\n"))
		}
	}
	w.Header().Set("Content-Type", "text/plain")
}

func (s *server) readUserTimelineHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := s.Logger(ctx)
	logger.Info("entering wkr2-api/user-timeline/read")

	params := validateReadTimelineParams(w, r)
	if params == nil {
		return
	}
	posts, err := s.userTimelineService.Get().ReadUserTimeline(ctx, params.reqID, params.userID, params.start, params.stop)
	if err != nil {
		http.Error(w, "error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// response, err := json.Marshal(posts)
	// iterate to add new line in between posts
	for i, post := range posts {
		postJSON, err := json.Marshal(post)
		if err != nil {
			http.Error(w, "error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(postJSON)
		if i < len(posts)-1 {
			w.Write([]byte("\n"))
		}
	}
	w.Header().Set("Content-Type", "text/plain")
}
