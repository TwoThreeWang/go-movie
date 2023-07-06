package zyw

import (
	"encoding/json"
	"fmt"
	"movie/configs"
	"movie/utils/easyhttp"
	"movie/utils/easylog"
	"time"
)

type ZYW struct{}

// ExternalGetReport 主动查询三方结果
func (zy *ZYW) ExternalGetReport(site configs.SpiderSites, kw string) (movie ZySearchResult, err error) {
	url := site.BaseUrl + "/provide/vod/?ac=videolist&pg=1&wd=" + kw
	param := map[string]interface{}{}
	header := map[string]string{
		"Content-Type": "application/json",
	}
	// Get请求接口
	jsonStr, err := easyhttp.Get(url, header, param)
	if err != nil {
		easylog.Log.Error(err)
		return
	}
	// 解析JSON响应
	err = json.Unmarshal(jsonStr, &movie)
	if err != nil {
		fmt.Println(url + " 错误结果 " + string(jsonStr))
		easylog.Log.Error(err)
		return
	}
	movie.Source = site.Key
	return
}

// ExternalGetById 根据影片ID查询三方数据
func (zy *ZYW) ExternalGetById(site configs.SpiderSites, vid string) (movie ZyDetailResult, err error) {
	url := site.BaseUrl + "/provide/vod/?ac=detail&ids=" + vid
	param := map[string]interface{}{}
	header := map[string]string{
		"Content-Type": "application/json",
	}
	// Get请求接口
	jsonStr, err := easyhttp.Get(url, header, param)
	if err != nil {
		easylog.Log.Error(err)
		return
	}
	// 解析JSON响应
	err = json.Unmarshal(jsonStr, &movie)
	if err != nil {
		fmt.Println(url + " 错误结果 " + string(jsonStr))
		easylog.Log.Error(err)
		return
	}
	movie.Source = site.Key
	return
}

type ZySearchResult struct {
	Code      int       `json:"code"`
	Msg       string    `json:"msg"`
	Page      string    `json:"page"`
	Pagecount int       `json:"pagecount"`
	Limit     string    `json:"limit"`
	Total     int       `json:"total"`
	Source    string    `json:"source"`
	List      []ZyMovie `json:"list"`
}
type ZyDetailResult struct {
	Code      int       `json:"code"`
	Msg       string    `json:"msg"`
	Page      int       `json:"page"`
	Pagecount int       `json:"pagecount"`
	Limit     string    `json:"limit"`
	Total     int       `json:"total"`
	Source    string    `json:"source"`
	List      []ZyMovie `json:"list"`
}
type ZyMovie struct {
	Source           string    `json:"source"`
	VodId            int       `json:"vod_id"`
	TypeId           int       `json:"type_id"`
	TypeId1          int       `json:"type_id_1"`
	GroupId          int       `json:"group_id"`
	VodName          string    `json:"vod_name"`
	VodSub           string    `json:"vod_sub"`
	VodEn            string    `json:"vod_en"`
	VodStatus        int       `json:"vod_status"`
	VodLetter        string    `json:"vod_letter"`
	VodColor         string    `json:"vod_color"`
	VodTag           string    `json:"vod_tag"`
	VodClass         string    `json:"vod_class"`
	VodPic           string    `json:"vod_pic"`
	VodPicThumb      string    `json:"vod_pic_thumb"`
	VodPicSlide      string    `json:"vod_pic_slide"`
	VodPicScreenshot string    `json:"vod_pic_screenshot"`
	VodActor         string    `json:"vod_actor"`
	VodDirector      string    `json:"vod_director"`
	VodWriter        string    `json:"vod_writer"`
	VodBehind        string    `json:"vod_behind"`
	VodBlurb         string    `json:"vod_blurb"`
	VodRemarks       string    `json:"vod_remarks"`
	VodPubdate       string    `json:"vod_pubdate"`
	VodTotal         int       `json:"vod_total"`
	VodSerial        string    `json:"vod_serial"`
	VodTv            string    `json:"vod_tv"`
	VodWeekday       string    `json:"vod_weekday"`
	VodArea          string    `json:"vod_area"`
	VodLang          string    `json:"vod_lang"`
	VodYear          string    `json:"vod_year"`
	VodVersion       string    `json:"vod_version"`
	VodState         string    `json:"vod_state"`
	VodAuthor        string    `json:"vod_author"`
	VodJumpurl       string    `json:"vod_jumpurl"`
	VodTpl           string    `json:"vod_tpl"`
	VodTplPlay       string    `json:"vod_tpl_play"`
	VodTplDown       string    `json:"vod_tpl_down"`
	VodIsend         int       `json:"vod_isend"`
	VodLock          int       `json:"vod_lock"`
	VodLevel         int       `json:"vod_level"`
	VodCopyright     int       `json:"vod_copyright"`
	VodPoints        int       `json:"vod_points"`
	VodPointsPlay    int       `json:"vod_points_play"`
	VodPointsDown    int       `json:"vod_points_down"`
	VodHits          int       `json:"vod_hits"`
	VodHitsDay       int       `json:"vod_hits_day"`
	VodHitsWeek      int       `json:"vod_hits_week"`
	VodHitsMonth     int       `json:"vod_hits_month"`
	VodDuration      string    `json:"vod_duration"`
	VodUp            int       `json:"vod_up"`
	VodDown          int       `json:"vod_down"`
	VodScore         string    `json:"vod_score"`
	VodScoreAll      int       `json:"vod_score_all"`
	VodScoreNum      int       `json:"vod_score_num"`
	VodTime          string    `json:"vod_time"`
	VodTimeAdd       int       `json:"vod_time_add"`
	VodTimeHits      int       `json:"vod_time_hits"`
	VodTimeMake      int       `json:"vod_time_make"`
	VodTrysee        int       `json:"vod_trysee"`
	VodDoubanId      int       `json:"vod_douban_id"`
	VodDoubanScore   string    `json:"vod_douban_score"`
	VodReurl         string    `json:"vod_reurl"`
	VodRelVod        string    `json:"vod_rel_vod"`
	VodRelArt        string    `json:"vod_rel_art"`
	VodPwd           string    `json:"vod_pwd"`
	VodPwdUrl        string    `json:"vod_pwd_url"`
	VodPwdPlay       string    `json:"vod_pwd_play"`
	VodPwdPlayUrl    string    `json:"vod_pwd_play_url"`
	VodPwdDown       string    `json:"vod_pwd_down"`
	VodPwdDownUrl    string    `json:"vod_pwd_down_url"`
	VodContent       string    `json:"vod_content"`
	VodPlayFrom      string    `json:"vod_play_from"`
	VodPlayServer    string    `json:"vod_play_server"`
	VodPlayNote      string    `json:"vod_play_note"`
	VodPlayUrl       string    `json:"vod_play_url"`
	VodDownFrom      string    `json:"vod_down_from"`
	VodDownServer    string    `json:"vod_down_server"`
	VodDownNote      string    `json:"vod_down_note"`
	VodDownUrl       string    `json:"vod_down_url"`
	VodPlot          int       `json:"vod_plot"`
	VodPlotName      string    `json:"vod_plot_name"`
	VodPlotDetail    string    `json:"vod_plot_detail"`
	TypeName         string    `json:"type_name"`
	UpdatedAt        time.Time `json:"updated_at"`
}
