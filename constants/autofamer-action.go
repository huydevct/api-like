package constants

type (
	// AutofarmerAction : định nghĩa các autofarmer action
	AutofarmerAction string

	// ExtraAutofarmerAction : định nghĩa các autofarmer action config manual
	ExtraAutofarmerAction string
)

// Autofarmer action
const (
	// Autolike
	//follow
	FollowAction         AutofarmerAction = "follow"
	FollowLowAction      AutofarmerAction = "follow_low"
	FollowHighAction     AutofarmerAction = "follow_high"
	FollowWarrantyAction AutofarmerAction = "follow_warranty"
	//likepage
	LikepageAction         AutofarmerAction = "likepage"
	LikepageLowAction      AutofarmerAction = "likepage_low"
	LikepageHighAction     AutofarmerAction = "likepage_high"
	LikePageWarrantyAction AutofarmerAction = "likepage_warranty"
	//
	BufflikeAction       AutofarmerAction = "bufflike"
	BuffcommentAction    AutofarmerAction = "buffcomment"
	ViplikeServiceAction AutofarmerAction = "viplikeService"

	// Autofarmer
	FeedlikeAction         AutofarmerAction = "FeedLike"
	WatchAction            AutofarmerAction = "Watch"
	WatchFollowAction      AutofarmerAction = "WatchFollow"
	WatchLikeAction        AutofarmerAction = "WatchLike"
	FeedAction             AutofarmerAction = "Feed"
	AddFriendSuggestAction AutofarmerAction = "AddFriendSuggest"
	ConfirmFriendAction    AutofarmerAction = "ConfirmFriend"
	PostStatusImageAction  AutofarmerAction = "PostStatusImage"
	PostStatusAction       AutofarmerAction = "PostStatus"
	// Page
	PageLikeAction AutofarmerAction = "PageLike"
	// Group
	GroupLikeAction AutofarmerAction = "GroupLike"
	JoinGroupAction AutofarmerAction = "JoinGroup"
	// Friend
	FriendLikeAction   AutofarmerAction = "FriendLike"
	AddFriendUIDAction AutofarmerAction = "AddFriendUID"
	// Share: page, group
	ShareNowAction   AutofarmerAction = "ShareNow"
	ShareGroupAction AutofarmerAction = "ShareGroup"

	//Action Insta
	FollowActionInsta     AutofarmerAction = "Follow"
	BuffLikeActionInsta   AutofarmerAction = "BuffLike"
	UnFollowActionInsta   AutofarmerAction = "UnFollow"
	CommentActionInsta    AutofarmerAction = "Comment"
	CommentTagActionInsta AutofarmerAction = "CommentTag"
	SearchActionInsta     AutofarmerAction = "Search"

	//Action Tiktok
	LikeActionTiktok      AutofarmerAction = "Like"
	CommentActionTiktok   AutofarmerAction = "Comment"
	FeedWatchActionTiktok AutofarmerAction = "FeedWatch"
)

// Extra action
const (
	ChangePasswordAction  ExtraAutofarmerAction = "ChangePassword"
	ChangeSecretkeyAction ExtraAutofarmerAction = "ChangeSecretkey"
	ChangeCoverAction     ExtraAutofarmerAction = "ChangeCover"
	ChangeAvatarAction    ExtraAutofarmerAction = "ChangeAvatar"
)

func (m AutofarmerAction) String() string {
	return string(m)
}

func (m ExtraAutofarmerAction) String() string {
	return string(m)
}
