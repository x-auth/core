package auth

//func Logout(w http.ResponseWriter, request *http.Request){
//	hydraAdmin := GetAdmin()
//	query := request.URL.Query()
//	// the challenge is used to fetch information from hydra
//	challenge_list := query["challenge"]
//	var challenge string
//	if len(challenge_list) == 0{
//		helpers.Error(w, 400, "Expected a login challenge but received none")
//		return
//	} else {
//		challenge = challenge_list[0]
//	}
//
//	logoutResponse, err := helpers.AcceptLogoutRequest(hydraAdmin, challenge)
//	if err != nil{
//		helpers.Error(w, 500, "Error logging in" + err.Error())
//		return
//	}
//	http.Redirect(w, request, *logoutResponse.GetPayload().RedirectTo, 307)
//}
