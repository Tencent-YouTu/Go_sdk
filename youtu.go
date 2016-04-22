/*
* File Name:	youtu.go
* Description:  http://open.youtu.qq.com API
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2015-06-19
 */

package youtu

import (
	"encoding/base64"
	"errors"
	"strconv"
)

const (
	//UserIDMaxLen 用户ID的最大长度
	UserIDMaxLen = 110
)

const expiredInterval = 1000

var (
	//ErrUserIDTooLong 用户ID过长错误
	ErrUserIDTooLong = errors.New("user id too long")
)

var (
	//DefaultHost 默认host
	DefaultHost = "http://api.youtu.qq.com"
)

var (
	//腾讯云host
	TencentYunHost = "https://youtu.api.qcloud.com"
)

//AppSign 应用签名鉴权
type AppSign struct {
	appID     uint32 //接入优图服务时,生成的唯一id, 用于唯一标识接入业务
	secretID  string //标识api鉴权调用者的密钥身份
	secretKey string //用于加密签名字符串和服务器端验证签名字符串的密钥，secret_key 必须严格保管避免泄露
	userID    string //接入业务自行定义的用户id，用于唯一标识一个用户, 登陆开发者账号的QQ号码
}

//NewAppSign 新建应用签名
func NewAppSign(appID uint32, secretID string, secretKey string, userID string) (as AppSign, err error) {
	if len(userID) > UserIDMaxLen {
		err = ErrUserIDTooLong
		return
	}
	as = AppSign{
		appID:     appID,
		secretID:  secretID,
		secretKey: secretKey,
		userID:    userID,
	}
	return
}

//Youtu 存储签名和host
type Youtu struct {
	appSign AppSign
	host    string
	debug   bool //Default false
}

func (y *Youtu) appID() string {
	return strconv.Itoa(int(y.appSign.appID))
}

//Init Youtu初始化
func Init(appSign AppSign, host string) *Youtu {
	return &Youtu{
		appSign: appSign,
		host:    host,
		debug:   false,
	}
}

//detectMode 检测模式，分正常和大脸
type detectMode int

const (
	//detectModeNormal 正常模式
	detectModeNormal detectMode = iota
	//detectModeBigFace 大脸模式
	detectModeBigFace
)

func mode(isBigFace bool) detectMode {
	if isBigFace {
		return detectModeBigFace
	}
	return detectModeNormal
}

// SetDebug For Debug
func (y *Youtu) SetDebug(isDebug bool) {
	y.debug = isDebug
}

type detectFaceReq struct {
	AppID string     `json:"app_id"`          //App的 API ID
	Image string     `json:"image,omitempty"` //base64编码的二进制图片数据
	Mode  detectMode `json:"mode,omitempty"`  //检测模式 0/1 正常/大脸模式
    Url   string     `json:"url,omitempty"`   //图片的url
}

//Face 脸参数
type Face struct {
	FaceID     string  `json:"face_id"`    //人脸标识
	X          int32   `json:"x"`          //人脸框左上角x
	Y          int32   `json:"y"`          //人脸框左上角y
	Width      float32 `json:"width"`      //人脸框宽度
	Height     float32 `json:"height"`     //人脸框高度
	Gender     int32   `json:"gender"`     //性别 [0/(female)~100(male)]
	Age        int32   `json:"age"`        //年龄 [0~100]
	Expression int32   `json:"expression"` //object 	微笑[0(normal)~50(smile)~100(laugh)]
	Glass      bool    `json:"glass"`      //是否有眼镜 [true,false]
	Pitch      int32   `json:"pitch"`      //上下偏移[-30,30]
	Yaw        int32   `json:"yaw"`        //左右偏移[-30,30]
	Roll       int32   `json:"roll"`       //平面旋转[-180,180]
    Beauty     int32   `json:"beauty"`     //魅力值 [0~100]
}

//DetectFaceRsp 脸检测返回
type DetectFaceRsp struct {
	SessionID   string `json:"session_id"`   //相应请求的session标识符，可用于结果查询
	ImageWidth  int32  `json:"image_width"`  //请求图片的宽度
	ImageHeight int32  `json:"image_height"` //请求图片的高度
	Face        []Face `json:"face"`         //被检测出的人脸Face的列表
	ErrorCode   int    `json:"errorcode"`    //返回状态值
	ErrorMsg    string `json:"errormsg"`     //返回错误消息
}

//DetectFace 检测给定图片(Image)中的所有人脸(Face)的位置和相应的面部属性。
//位置包括(x, y, w, h)，面部属性包括性别(gender), 年龄(age), 魅力值(beauty)
//表情(expression), 眼镜(glass)和姿态(pitch，roll，yaw).
//imageType 表示image类型是图片还是URL, 其中0代表图片,1代表url
func (y *Youtu) DetectFace(image []byte, isBigFace bool, imageType int) (rsp DetectFaceRsp, err error) {
	var req detectFaceReq
    req.AppID = strconv.Itoa(int(y.appSign.appID))
    req.Mode = mode(isBigFace)
    
    if imageType == 0 {
        req.Image = base64.StdEncoding.EncodeToString(image)
    } else {
        req.Url = string(image)
    }
    
    err = y.interfaceRequest("detectface", req, &rsp, 0)
	return
}

type faceShapeReq struct {
	AppID string     `json:"app_id"`          //App的 API ID
	Image string     `json:"image,omitempty"` //base64编码的二进制图片数据
	Mode  detectMode `json:"mode,omitempty"`  //检测模式 0/1 正常/大脸模式
    Url   string     `json:"url,omitempty"`   //图片的url
}

type pos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

//FaceShape 五官定位
type FaceShape struct {
	FaceProfile  []pos `json:"face_profile"`  //描述脸型轮廓的21点
	LeftEye      []pos `json:"left_eye"`      //描述左眼轮廓的8点
	RightEye     []pos `json:"right_eye"`     //描述右眼轮廓的8点
	LeftEyebrow  []pos `json:"left_eyebrow"`  //描述左眉轮廓的8点
	RightEyebrow []pos `json:"right_eyebrow"` //描述右眉轮廓的8点
	Mouth        []pos `json:"mouth"`         //描述嘴巴轮廓的22点
	Nose         []pos `json:"nose"`          //描述鼻子轮廓的13点
}

// FaceShapeRsp 返回
type FaceShapeRsp struct {
	SessionID   string      `json:"session_id"`   //相应请求的session标识符，可用于结果查询
	FaceShape   []FaceShape `json:"face_shape"`   //人脸轮廓结构体，包含所有人脸的轮廓点
	ImageWidth  int         `json:"image_width"`  //请求图片的宽度
	ImageHeight int         `json:"image_height"` //请求图片的高度
	ErrorCode   int         `json:"errorcode"`    //返回状态值
	ErrorMsg    string      `json:"errormsg"`     //返回错误消息
}

//FaceShape 对请求图片进行五官定位，计算构成人脸轮廓的88个点，包括眉毛（左右各8点）、眼睛（左右各8点）、鼻子（13点）、嘴巴（22点）、脸型轮廓（21点）
//imageType 表示image的类型是图片还是URL, 其中0代表图片,1代表url
func (y *Youtu) FaceShape(image []byte, isBigFace bool, imageType int) (rsp FaceShapeRsp, err error) {
	var req faceShapeReq
    req.AppID = strconv.Itoa(int(y.appSign.appID))
    req.Mode = mode(isBigFace)
    
    if imageType == 0 {
        req.Image = base64.StdEncoding.EncodeToString(image)
    } else {
        req.Url =  string(image)
    }
    
	err = y.interfaceRequest("faceshape", req, &rsp, 0)
	return
}

type faceCompareReq struct {
	AppID  string `json:"app_id"`
	ImageA string `json:"imageA,omitempty"` //使用base64编码的二进制图片数据A
	ImageB string `json:"imageB,omitempty"` //使用base64编码的二进制图片数据B
    UrlA   string `json:"urlA,omitempty"`   //图片A的url
    UrlB   string `json:"urlB,omitempty"`   //图片B的url
}

//FaceCompareRsp 脸比较返回
type FaceCompareRsp struct {
	SessionID  string  `json:"session_id"`   //相应请求的session标识符，可用于结果查询
    Similarity float32 `json:"similarity"`  //两个face的相似度
	ErrorCode  int32   `json:"errorcode"`   //返回状态码
	ErrorMsg   string  `json:"errormsg"`    //返回错误消息
}

//FaceCompare 计算两个Face的相似性以及五官相似度
//imageType 表示image类型是图片还是URL, 其中0代表图片,1代表url
func (y *Youtu) FaceCompare(imageA, imageB []byte, imageType int) (rsp FaceCompareRsp, err error) {
	var req faceCompareReq
    req.AppID = y.appID()
   
    if imageType == 0 {
        req.ImageA = base64.StdEncoding.EncodeToString(imageA)
        req.ImageB = base64.StdEncoding.EncodeToString(imageB)
    } else {
        req.UrlA = string(imageA)
        req.UrlB = string(imageB)
    }
	
    err = y.interfaceRequest("facecompare", req, &rsp, 0)
	return
}

type faceVerifyReq struct {
	AppID    string `json:"app_id"`    //App的 API ID
	Image    string `json:"image"`     //使用base64编码的二进制图片数据
	PersonID string `json:"person_id"` //待验证的Person
    Url   string    `json:"url,omitempty"`   //图片的url
}

//FaceVerifyRsp 脸验证返回
type FaceVerifyRsp struct {
	Ismatch    bool    `json:"ismatch"`    //两个输入是否为同一人的判断
	Confidence float32 `json:"confidence"` //系统对这个判断的置信度。
	SessionID  string  `json:"session_id"` //相应请求的session标识符，可用于结果查询
	ErrorCode  int32   `json:"errorcode"`  //返回状态码
	ErrorMsg   string  `json:"errormsg"`   //返回错误消息
}

//FaceVerify 给定一个Face和一个Person，返回是否是同一个人的判断以及信度。
//imageType 表示image类型是图片还是URL, 其中0代表图片,1代表url
func (y *Youtu) FaceVerify(personID string, image []byte, imageType int) (rsp FaceVerifyRsp, err error) {
	var req faceVerifyReq
    req.AppID = y.appID()
    req.PersonID = personID

    if  imageType == 0 {
        req.Image = base64.StdEncoding.EncodeToString(image)
    } else {
        req.Url = string(image)
    }
	
    err = y.interfaceRequest("faceverify", req, &rsp, 0)
	return
}

type faceIdentifyReq struct {
	AppID   string `json:"app_id"`   //App的 API ID
	GroupID string `json:"group_id"` //候选人组id
	Image   string `json:"image,omitempty"`    //使用base64编码的二进制图片数据
    Url   string   `json:"url,omitempty"`   //图片的url
}

type Candidate struct {
    PersonID   string `json:"person_id"`  //识别结果，person_id
    FaceID     string  `json:"face_id"`    //识别的face_id
    Confidence float32 `json:"confidence"` //置信度
    Tag        string  `json:"tag"` //人脸备注信息[]
}
//FaceIdentifyRsp 脸识别返回
type FaceIdentifyRsp struct {
	SessionID  string  `json:"session_id"` //相应请求的session标识符，可用于结果查询
	Candidates []Candidate `json:"candidates"` //识别出的top5候选人
    ErrorCode  int     `json:"errorcode"`  //返回状态码
	ErrorMsg   string  `json:"errormsg"`   //返回错误消息
}

//FaceIdentify 对于一个待识别的人脸图片，在一个Group中识别出最相似的Person作为其身份返回
//imageType 表示image类型是图片还是URL, 其中0代表图片,1代表url
func (y *Youtu) FaceIdentify(groupID string, image []byte, imageType int) (rsp FaceIdentifyRsp, err error) {
	var req faceIdentifyReq
    req.AppID = y.appID()
    req.GroupID = groupID

    if imageType == 0 {
        req.Image = base64.StdEncoding.EncodeToString(image)
    } else {
        req.Url = string(image)
    }
	
    err = y.interfaceRequest("faceidentify", req, &rsp, 0)
	return
}

type newPersonReq struct {
	AppID      string   `json:"app_id"` //App的 API ID
	Image      string   `json:"image,omitempty"`  //使用base64编码的二进制图片数据
	PersonID   string   `json:"person_id"`
	GroupIDs   []string `json:"group_ids"`             // 	加入到组的列表
	PersonName string   `json:"person_name,omitempty"` //名字
	Tag        string   `json:"tag,omitempty"`         //备注信息
    Url        string   `json:"url,omitempty"`   //图片的url
}

//NewPersonRsp 个体创建返回
type NewPersonRsp struct {
	SessionID  string `json:"session_id"`  //相应请求的session标识符
	SucGroup   int    `json:"suc_group"`   //成功被加入的group数量
	SucFace    int    `json:"suc_face"`    //成功加入的face数量
	PersonID   string `json:"person_id"`   //相应person的id
	FaceID     string `json:"face_id"`     //创建所用图片生成的face_id
	GroupIds   []string `json:"group_ids"`   //加入成功的组id
	ErrorCode  int    `json:"errorcode"`   //返回码
	ErrorMsg   string `json:"errormsg"`    //返回错误消息
}

//NewPerson 创建一个Person，并将Person放置到group_ids指定的组当中
//imageType 表示image类型是图片还是URL, 其中0代表图片,1代表url
func (y *Youtu) NewPerson(personID string, personName string, groupIDs []string, image []byte, tag string, imageType int) (rsp NewPersonRsp, err error) {
	var req newPersonReq
    req.AppID = y.appID()
    req.PersonID = personID
    req.GroupIDs = groupIDs
    req.PersonName = personName
    req.Tag = tag

    if imageType == 0 {
        req.Image = base64.StdEncoding.EncodeToString(image)
    } else {
        req.Url = string(image)
    }
	
    err = y.interfaceRequest("newperson", req, &rsp, 0)
	return
}

type delPersonReq struct {
	AppID    string `json:"app_id"`
	PersonID string `json:"person_id"` //待删除个体ID
}

//DelPersonRsp 删除个体返回
type DelPersonRsp struct {
	SessionID string `json:"session_id"` //相应请求的session标识符
	Deleted   int    `json:"deleted"`    //成功删除的Person数量
    PersonID   string `json:"person_id"`   //相应person的id
    ErrorCode int    `json:"errorcode"`  //返回状态码
	ErrorMsg  string `json:"errormsg"`   //返回错误消息
}

//DelPerson 删除一个Person
func (y *Youtu) DelPerson(personID string) (rsp DelPersonRsp, err error) {
	req := delPersonReq{
		AppID:    y.appID(),
		PersonID: personID,
	}
	err = y.interfaceRequest("delperson", req, &rsp, 0)
	return
}

type addFaceReq struct {
	AppID    string   `json:"app_id"`        //App的 API ID
	PersonID string   `json:"person_id"`     //String 	待增加人脸的个体id
	Images   []string `json:"images,omitempty"`        //base64编码的二进制图片数据构成的数组
	Tag      string   `json:"tag,omitempty"` //备注信息
    Urls      []string `json:"urls,omitempty"`   //图片的url
}

//AddFaceRsp 增加人脸返回
type AddFaceRsp struct {
	SessionID string   `json:"session_id"` //相应请求的session标识符
	Added     int      `json:"added"`      //成功加入的face数量
	FaceIDs   []string `json:"face_ids"`   //增加的人脸ID列表
	RetCodes  []int    `json:"ret_codes"`  //每张图片增加人脸的返回码[]
    ErrorCode int      `json:"errorcode"`  //返回状态码
	ErrorMsg  string   `json:"errormsg"`   //返回错误消息
}

//AddFace 将一组Face加入到一个Person中。注意，一个Face只能被加入到一个Person中。
//一个Person最多允许包含10000个Face
//imageType 表示image类型是图片还是URL, 其中0代表图片,1代表url
func (y *Youtu) AddFace(personID string, images [][]byte, tag string, imageType int) (rsp AddFaceRsp, err error) {
	var req addFaceReq
    req.AppID =  y.appID()
    req.PersonID = personID
    req.Tag = tag
    
    imageDatas := make([]string, len(images))
    if imageType == 0 {
        for i, img := range images {
            imageDatas[i] = base64.StdEncoding.EncodeToString([]byte(img))
        }
        req.Images = imageDatas
   
   } else {
        for i, img := range images {
            imageDatas[i] = string([]byte(img))
        }
        req.Urls = imageDatas
    }
	
    err = y.interfaceRequest("addface", req, &rsp, 0)
	return
}

type delFaceReq struct {
	AppID    string   `json:"app_id"`    //App的 API ID
	PersonID string   `json:"person_id"` //待删除人脸的person ID
	FaceIDs  []string `json:"face_ids"`  //删除人脸id的列表
}

//DelFaceRsp 删除人脸返回
type DelFaceRsp struct {
	SessonID  string `json:"session_id"` //相应请求的session标识符
	Deleted   int32  `json:"deleted"`    //成功删除的face数量
    FaceIDs   []string `json:"face_ids"` //成功删除的人脸ID列表
    ErrorCode int32  `json:"errorcode"`  //返回状态码
	ErrorMsg  string `json:"errormsg"`   //返回错误消息
}

//DelFace 删除一个person下的face，包括特征，属性和face_id.
func (y *Youtu) DelFace(personID string, faceIDs []string) (rsp DelFaceRsp, err error) {
	req := delFaceReq{
		AppID:    y.appID(),
		PersonID: personID,
		FaceIDs:  faceIDs,
	}
	err = y.interfaceRequest("delface", req, &rsp, 0)
	return
}

type setInfoReq struct {
	AppID      string `json:"app_id"` //App的 API ID
	PersonID   string `json:"person_id"`
	PersonName string `json:"person_name,omitempty"` //新的name
	Tag        string `json:"tag,omitempty"`         //备注信息
}

//SetInfoRsp 设置信息返回
type SetInfoRsp struct {
    SessonID  string `json:"session_id"` //相应请求的session标识符
    PersonID string   `json:"person_id"` //待删除人脸的person ID
    ErrorCode int32  `json:"errorcode"`  //返回状态码
    ErrorMsg  string `json:"errormsg"`   //返回错误消息
}

//SetInfo 设置Person的name.
func (y *Youtu) SetInfo(personID string, personName string, tag string) (rsp SetInfoRsp, err error) {
	req := setInfoReq{
		AppID:      y.appID(),
		PersonID:   personID,
		PersonName: personName,
		Tag:        tag,
	}
	err = y.interfaceRequest("setinfo", req, &rsp, 0)
	return
}

type getInfoReq struct {
	AppID    string `json:"app_id"`    //App的 API ID
	PersonID string `json:"person_id"` //待查询个体的ID
}

//GetInfoRsp 获取信息返回
type GetInfoRsp struct {
	PersonName string   `json:"person_name"` //相应person的name
	PersonID   string   `json:"person_id"`   //相应person的id
	GroupIDs   []string `json:"group_ids"`   //包含此个体的组列表
	FaceIDs    []string `json:"face_ids"`    //包含的人脸列表
	SessionID  string   `json:"session_id"` //相应请求的session标识符
	ErrorCode  int      `json:"errorcode"` //返回状态码
	ErrorMsg   string   `json:"errormsg"`  //返回错误消息
}

//GetInfo 获取一个Person的信息, 包括name, id, tag, 相关的face, 以及groups等信息。
func (y *Youtu) GetInfo(personID string) (rsp GetInfoRsp, err error) {
	req := getInfoReq{
		AppID:    y.appID(),
		PersonID: personID,
	}
	err = y.interfaceRequest("getinfo", req, &rsp, 0)
	return
}

type getGroupIDsReq struct {
	AppID string `json:"app_id"` //App的 API ID
}

//GetGroupIDsRsp 获取组ID返回
type GetGroupIDsRsp struct {
	GroupIDs  []string `json:"group_ids"` //相应app_id的group_id列表
	ErrorCode int32    `json:"errorcode"` //返回状态码
	ErrorMsg  string   `json:"errormsg"`  //返回错误消息
}

//GetGroupIDs 获取一个appId下所有group列表
func (y *Youtu) GetGroupIDs() (rsp GetGroupIDsRsp, err error) {
	req := getGroupIDsReq{
		AppID: y.appID(),
	}
	err = y.interfaceRequest("getgroupids", req, &rsp, 0)
	return
}

type getPersonIDsReq struct {
	AppID   string `json:"app_id"`   //App的 API ID
	GroupID string `json:"group_id"` //组id
}

//GetPersonIDsRsp 获取个人ID返回
type GetPersonIDsRsp struct {
	PersonIDs []string `json:"person_ids"` //相应person的id列表
	ErrorCode int32    `json:"errorcode"`  //返回状态码
	ErrorMsg  string   `json:"errormsg"`   //返回错误消息
}

//GetPersonIDs 获取一个组Group中所有person列表
func (y *Youtu) GetPersonIDs(groupID string) (rsp GetPersonIDsRsp, err error) {
	req := getPersonIDsReq{
		AppID:   y.appID(),
		GroupID: groupID,
	}
	err = y.interfaceRequest("getpersonids", req, &rsp, 0)
	return
}

type getFaceIDsReq struct {
	AppID    string `json:"app_id"`    //App的 API ID
	PersonID string `json:"person_id"` //个体id
}

//GetFaceIDsRsp 获取脸ID返回
type GetFaceIDsRsp struct {
	FaceIDs   []string `json:"face_ids"`  //相应face的id列表
	ErrorCode int32    `json:"errorcode"` //返回状态码
	ErrorMsg  string   `json:"errormsg"`  //返回错误消息
}

//GetFaceIDs 获取一个组person中所有face列表
func (y *Youtu) GetFaceIDs(personID string) (rsp GetFaceIDsRsp, err error) {
	req := getFaceIDsReq{
		AppID:    y.appID(),
		PersonID: personID,
	}
	err = y.interfaceRequest("getfaceids", req, &rsp, 0)
	return
}

type getFaceInfoReq struct {
	AppID  string `json:"app_id"`  //App的 API ID
	FaceID string `json:"face_id"` //人脸id
}

//GetFaceInfoRsp 获取脸部信息返回
type GetFaceInfoRsp struct {
	FaceInfo  Face   `json:"face_info"` //人脸信息
	ErrorCode int32  `json:"errorcode"` //返回状态码
	ErrorMsg  string `json:"errormsg"`  //返回错误消息
}

//GetFaceInfo 获取一个face的相关特征信息
func (y *Youtu) GetFaceInfo(faceID string) (rsp GetFaceInfoRsp, err error) {
	req := getFaceInfoReq{
		AppID:  y.appID(),
		FaceID: faceID,
	}
	err = y.interfaceRequest("getfaceinfo", req, &rsp, 0)
	return
}


type FuzzyDetectReq struct{
    AppID  string `json:"app_id"`  //App的 API ID
    Url    string `json:"url,omitempty"`   //图片的url
    Image  string `json:"image,omitempty"`  //使用base64编码的二进制图片数据
    Seq    string `json:"seq,omitempty"`  // 序列号
}

type FuzzyDetectRsp struct{
    Fuzzy bool `json:"fuzzy"`  // 是否模糊 
    FuzzyConfidence  float32 `json:"fuzzy_confidence"` //范围 0-1的浮点数,越大置信度越高
    ErrorCode int32  `json:"errorcode"` //返回状态码
    ErrorMsg  string `json:"errormsg"`  //返回错误消息
}
//FuzzyDetect 检测图片的模糊度
//imageType 表示image类型是图片还是URL, 其中0代表图片,1代表url
func (y *Youtu) FuzzyDetect(image []byte, imageType int, seq string) (rsp FuzzyDetectRsp, err error) {
    var req FuzzyDetectReq
    req.AppID =  y.appID()
    req.Seq = seq

    if imageType == 0 {
        req.Image = base64.StdEncoding.EncodeToString(image)
    } else {
        req.Url = string(image)
    }
    err = y.interfaceRequest("fuzzydetect", req, &rsp, 1)
    return 
}

type FoodDetectReq struct{
    AppID  string `json:"app_id"`  //App的 API ID
    Url    string `json:"url,omitempty"`   //图片的url
    Image  string `json:"image,omitempty"`  //使用base64编码的二进制图片数据
    Seq    string `json:"seq,omitempty"`  // 序列号
}

type FoodDetectRsp struct{
    Food bool `json:"food"`  // 是否美食 
    FoodConfidence  float32 `json:"food_confidence"` 
    ErrorCode int32  `json:"errorcode"` //返回状态码
    ErrorMsg  string `json:"errormsg"`  //返回错误消息
}

//FoodDetect 美食检测 
//imageType 表示image类型是图片还是URL, 其中0代表图片,1代表url
func (y *Youtu) FoodDetect(image []byte, imageType int, seq string) (rsp FoodDetectRsp, err error) {
    var req FoodDetectReq
    req.AppID =  y.appID()
    req.Seq = seq

    if imageType == 0 {
            req.Image = base64.StdEncoding.EncodeToString(image)
    } else {
            req.Url = string(image)
    }
    err = y.interfaceRequest("fooddetect", req, &rsp, 1)
    return
}

type ImageTagReq struct{
    AppID  string `json:"app_id"`  //App的 API ID
    Url    string `json:"url,omitempty"`   //图片的url 
    Image  string `json:"image,omitempty"`  //使用base64编码的二进制图片数据
    Seq    string `json:"seq,omitempty"`  // 序列号
}

type ImageTag struct {
    TagName string `json:"tag_name"`
    TagConfidence int `json:"tag_confidence"`
}

type ImageTagRsp struct{
    Seq    string `json:"seq,omitempty"`  // 序列号
    Tags   []ImageTag `json:"tags"`
    ErrorCode int32  `json:"errorcode"` //返回状态码
    ErrorMsg  string `json:"errormsg"`  //返回错误消息
}


//ImageTag 图片分类 
//imageType 表示image类型是图片还是URL, 其中0代表图片,1代表url
func (y *Youtu) ImageTag(image []byte, imageType int, seq string) (rsp ImageTagRsp, err error) {
    var req ImageTagReq
    req.AppID =  y.appID()
    req.Seq = seq 

    if imageType == 0 { 
       req.Image = base64.StdEncoding.EncodeToString(image)
    } else {
       req.Url = string(image)
    }   
    
    err = y.interfaceRequest("imagetag", req, &rsp, 1)
    return
}

type ImagePornReq struct{
    AppID  string `json:"app_id"`  //App的 API ID
    Url    string `json:"url,omitempty"`   //图片的url 
    Image  string `json:"image,omitempty"`  //使用base64编码的二进制图片数据
    Seq    string `json:"seq,omitempty"`  // 序列号
}


type ImagePornRsp struct{
    Seq    string `json:"seq,omitempty"`  // 序列号
    Tags   []ImageTag `json:"tags"`
    ErrorCode int32  `json:"errorcode"` //返回状态码
    ErrorMsg  string `json:"errormsg"`  //返回错误消息
}

//ImagePorn 图片分类 
//imageType 表示image类型是图片还是URL, 其中0代表图片,1代表url
func (y *Youtu) ImagePorn(image []byte, imageType int, seq string) (rsp ImagePornRsp, err error) {
    var req ImagePornReq
    req.AppID =  y.appID()
    req.Seq = seq 

    if imageType == 0 { 
       req.Image = base64.StdEncoding.EncodeToString(image)
    } else {
       req.Url = string(image)
    }   
    
    err = y.interfaceRequest("imageporn", req, &rsp, 1)
    return
}

type IdcardOcrReq struct{
    AppID  string `json:"app_id"`  //App的 API ID
    Url    string `json:"url,omitempty"`   //图片的url 
    Image  string `json:"image,omitempty"`  //使用base64编码的二进制图片数据
    CardType int32 `json:"card_type,omitempty"`        // 身份证正反面
    SessionId  string `json:"session_id,omitempty"`  // 序列号
} 

type IdcardOcrRsp struct{
    SessionId    string `json:"session_id,omitempty"`  // 序列号
    Name   string `json:"name,omitempty"`  // 证件姓名
    NameConfidenceAll []int32 `json:"name_confidence_all,omitempty"`  // 证件姓名置信度   
    Sex   string `json:"sex,omitempty"`  // 性别
    SexConfidenceAll []int32 `json:"sex_confidence_all,omitempty"`  // 性别置信度    
    Nation   string `json:"nation,omitempty"`  // 民族
    NationConfidenceAll []int32 `json:"nation_confidence_all,omitempty"`  // 民族置信度
    Birth   string `json:"birth,omitempty"`  // 出生日期
    BirthConfidenceAll []int32 `json:"birth_confidence_all,omitempty"`  // 出生日期置信度
    Adress   string `json:"address,omitempty"`  // 出生日期
    AddressConfidenceAll []int32 `json:"address_confidence_all,omitempty"`  // 出生日期置信度
    Id   string `json:"id,omitempty"`  // 身份证号
    IdConfidenceAll []int32 `json:"id_confidence_all,omitempty"`  // 身份证号置信度
    FrontImage   string `json:"frontimage,omitempty"`  // OCR识别的身份证正面照片
    FrontImageConfidenceAll []int32 `json:"frontimage_confidence_all,omitempty"`  // 正面照片置信度
    WaterMask   int32 `json:"watermask_status,omitempty"`  //水印是否存在(暂时不提供)
    WaterMaskConfidenceAll []int32 `json:"watermask_confidence_all,omitempty"`  // 水印置信度
    ValidDate   string `json:"valid_date,omitempty"`  //证件的有效期
    ValidDateConfidenceAll []int32 `json:"valid_date_confidence_all,omitempty"`  // 证件的有效期置信度
    Authority   string `json:"authority,omitempty"`  //发证机关
    AuthorityConfidenceAll []int32 `json:"authority_confidence_all,omitempty"`  // 发证机关置信度
    BackImage   string `json:"backimage,omitempty"`  //OCR识别的证件身份证反面照片
    BackImageConfidenceAll []int32 `json:"backimage_confidence_all,omitempty"`  // 反面照片置信度
    DetailErrorcode []int32 `json:"detail_errorcode,omitempty"`  //详细的错误原因
    DetailErrormsg []string `json:"detail_errormsg,omitempty"`  // 详细的错误原因说明
    ErrorCode int32  `json:"errorcode"` //返回状态码
    ErrorMsg  string `json:"errormsg"`  //返回错误消息
}

//IdcardOcr 图片分类 
//imageType 表示image类型是图片还是URL, 其中0代表图片,1代表url
//cardType 代表身份证正面还是反面，其中0代表证明，1代表反面
func (y *Youtu) IdcardOcr(image []byte, imageType int, cardType int32, seq string) (rsp IdcardOcrRsp, err error) {
    var req IdcardOcrReq
    req.AppID =  y.appID()
    req.SessionId = seq 
    req.CardType = cardType
    
    if imageType == 0 { 
       req.Image = base64.StdEncoding.EncodeToString(image)
    } else {
       req.Url = string(image)
    }   
    
    err = y.interfaceRequest("idcardocr", req, &rsp, 2)
    return
}

type NamecardOcrReq struct{
    AppID  string `json:"app_id"`  //App的 API ID
    Url    string `json:"url,omitempty"`   //图片的url 
    Image  string `json:"image,omitempty"`  //使用base64编码的二进制图片数据
    RetImage bool `json:"retimage,omitempty"`        // 是否需要返回处理结果图 
    SessionId  string `json:"session_id,omitempty"`  // 序列号
} 

type NamecardOcrRsp struct{
    SessionId    string `json:"session_id,omitempty"`  // 序列号
    Phone   string `json:"phone,omitempty"`  // 手机号
    PhoneConfidence float32 `json:"phone_confidence ,omitempty"`  // 手机号置信度   
    Name   string `json:"name ,omitempty"`  // 姓名
    NameConfidence float32 `json:"name_confidence ,omitempty"`  // 姓名置信度  
    Image    string `json:"image ,omitempty"`  // 处理结果图片
    ErrorCode int32  `json:"errorcode"` //返回状态码
    ErrorMsg  string `json:"errormsg"`  //返回错误消息
}

//NameCardOcr 图片分类 
//imageType 表示image类型是图片还是URL, 其中0代表图片,1代表url
//是否需要返回处理结果图,true 返回，false 不返回
func (y *Youtu) NameCardOcr(image []byte, imageType int, retImage bool, seq string) (rsp NamecardOcrRsp, err error) {
    var req NamecardOcrReq
    req.AppID =  y.appID()
    req.SessionId = seq 
    req.RetImage = retImage
    
    if imageType == 0 { 
       req.Image = base64.StdEncoding.EncodeToString(image)
    } else {
       req.Url = string(image)
    }   
    
    err = y.interfaceRequest("namecardocr", req, &rsp, 2)
    return
}
