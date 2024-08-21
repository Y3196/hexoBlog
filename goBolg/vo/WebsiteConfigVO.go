package vo

// WebsiteConfigVO represents the configuration of the website.
type WebsiteConfigVO struct {
	WebsiteAvatar     string   `json:"websiteAvatar"`     // Website Avatar
	WebsiteName       string   `json:"websiteName"`       // Website Name
	WebsiteAuthor     string   `json:"websiteAuthor"`     // Website Author
	WebsiteIntro      string   `json:"websiteIntro"`      // Website Introduction
	WebsiteNotice     string   `json:"websiteNotice"`     // Website Notice
	WebsiteCreateTime string   `json:"websiteCreateTime"` // Website Creation Time (String format, e.g., "2019-09-12T00:00:00Z")
	WebsiteRecordNo   string   `json:"websiteRecordNo"`   // Website Record Number
	SocialLoginList   []string `json:"socialLoginList"`   // Social Login List
	SocialUrlList     []string `json:"socialUrlList"`     // Social URL List
	QQ                string   `json:"qq"`                // QQ
	Github            string   `json:"github"`            // GitHub
	Gitee             string   `json:"gitee"`             // Gitee
	TouristAvatar     string   `json:"touristAvatar"`     // Tourist Avatar
	UserAvatar        string   `json:"userAvatar"`        // User Avatar
	IsCommentReview   int      `json:"isCommentReview"`   // Is Comment Review
	IsMessageReview   int      `json:"isMessageReview"`   // Is Message Review
	IsEmailNotice     int      `json:"isEmailNotice"`     // Is Email Notice
	IsReward          int      `json:"isReward"`          // Is Reward
	WeiXinQRCode      string   `json:"weiXinQRCode"`      // WeChat QR Code
	AlipayQRCode      string   `json:"alipayQRCode"`      // Alipay QR Code
	ArticleCover      string   `json:"articleCover"`      // Article Cover
	IsChatRoom        int      `json:"isChatRoom"`        // Is Chat Room
	WebsocketUrl      string   `json:"websocketUrl"`      // WebSocket URL
	IsMusicPlayer     int      `json:"isMusicPlayer"`     // Is Music Player
}
