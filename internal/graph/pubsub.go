package graph

import (
	"time"

	"github.com/moby/pubsub"
)

var (
	// UserUpdated                chan *model.User
	// UserDeleted                chan *model.User
	// TeamMemberAdded            chan *model.TeamMember
	// TeamMemberUpdated          chan *model.TeamMember
	// TeamMemberRemoved          chan string
	// TeamCollectionAdded        chan *model.TeamCollection
	// TeamCollectionUpdated      chan *model.TeamCollection
	// TeamCollectionRemoved      chan string
	// TeamCollectionMoved        chan *model.TeamCollection
	// CollectionOrderUpdated     chan *dto.CollectionReorderData
	// TeamRequestAdded           chan *model.TeamRequest
	// TeamRequestUpdated         chan *model.TeamRequest
	// TeamRequestDeleted         chan string
	// RequestOrderUpdated        chan *dto.RequestReorderData
	// RequestMoved               chan *model.TeamRequest
	// TeamInvitationAdded        chan *model.TeamInvitation
	// TeamInvitationRemoved      chan string
	// TeamEnvironmentUpdated     chan *model.TeamEnvironment
	// TeamEnvironmentCreated     chan *model.TeamEnvironment
	// TeamEnvironmentDeleted     chan *model.TeamEnvironment
	// UserSettingsCreated        chan *model.UserSetting
	// UserSettingsUpdated        chan *model.UserSetting
	// UserEnvironmentCreated     chan *model.UserEnvironment
	// UserEnvironmentUpdated     chan *model.UserEnvironment
	// UserEnvironmentDeleted     chan *model.UserEnvironment
	// UserEnvironmentDeleteMany  chan int
	// UserHistoryCreated         chan *model.UserHistory
	// UserHistoryUpdated         chan *model.UserHistory
	// UserHistoryDeleted         chan *model.UserHistory
	// UserHistoryDeletedMany     chan *dto.UserHistoryDeletedManyData
	// UserRequestCreated         chan *model.UserRequest
	// UserRequestUpdated         chan *model.UserRequest
	// UserRequestDeleted         chan *model.UserRequest
	// UserRequestMoved           chan *dto.UserRequestReorderData
	// UserCollectionCreated      chan *model.UserCollection
	// UserCollectionUpdated      chan *model.UserCollection
	// UserCollectionRemoved      chan *dto.UserCollectionRemovedData
	// UserCollectionMoved        chan *model.UserCollection
	// UserCollectionOrderUpdated chan *dto.UserCollectionReorderData
	// MyShortcodesCreated        chan *model.Shortcode
	// MyShortcodesUpdated        chan *model.Shortcode
	// MyShortcodesRevoked        chan *model.Shortcode

	UserUpdatedSub                *pubsub.Publisher
	UserDeletedSub                *pubsub.Publisher
	UserInvitedSub                *pubsub.Publisher
	TeamMemberAddedSub            *pubsub.Publisher
	TeamMemberUpdatedSub          *pubsub.Publisher
	TeamMemberRemovedSub          *pubsub.Publisher
	TeamCollectionAddedSub        *pubsub.Publisher
	TeamCollectionUpdatedSub      *pubsub.Publisher
	TeamCollectionRemovedSub      *pubsub.Publisher
	TeamCollectionMovedSub        *pubsub.Publisher
	CollectionOrderUpdatedSub     *pubsub.Publisher
	TeamRequestAddedSub           *pubsub.Publisher
	TeamRequestUpdatedSub         *pubsub.Publisher
	TeamRequestDeletedSub         *pubsub.Publisher
	RequestOrderUpdatedSub        *pubsub.Publisher
	RequestMovedSub               *pubsub.Publisher
	TeamInvitationAddedSub        *pubsub.Publisher
	TeamInvitationRemovedSub      *pubsub.Publisher
	TeamEnvironmentUpdatedSub     *pubsub.Publisher
	TeamEnvironmentCreatedSub     *pubsub.Publisher
	TeamEnvironmentDeletedSub     *pubsub.Publisher
	UserSettingsCreatedSub        *pubsub.Publisher
	UserSettingsUpdatedSub        *pubsub.Publisher
	UserEnvironmentCreatedSub     *pubsub.Publisher
	UserEnvironmentUpdatedSub     *pubsub.Publisher
	UserEnvironmentDeletedSub     *pubsub.Publisher
	UserEnvironmentDeleteManySub  *pubsub.Publisher
	UserHistoryCreatedSub         *pubsub.Publisher
	UserHistoryUpdatedSub         *pubsub.Publisher
	UserHistoryDeletedSub         *pubsub.Publisher
	UserHistoryDeletedManySub     *pubsub.Publisher
	UserRequestCreatedSub         *pubsub.Publisher
	UserRequestUpdatedSub         *pubsub.Publisher
	UserRequestDeletedSub         *pubsub.Publisher
	UserRequestMovedSub           *pubsub.Publisher
	UserCollectionCreatedSub      *pubsub.Publisher
	UserCollectionUpdatedSub      *pubsub.Publisher
	UserCollectionRemovedSub      *pubsub.Publisher
	UserCollectionMovedSub        *pubsub.Publisher
	UserCollectionOrderUpdatedSub *pubsub.Publisher
	MyShortcodesCreatedSub        *pubsub.Publisher
	MyShortcodesUpdatedSub        *pubsub.Publisher
	MyShortcodesRevokedSub        *pubsub.Publisher
)

func init() {
	// UserUpdated = make(chan *model.User)
	// UserDeleted = make(chan *model.User)
	// TeamMemberAdded = make(chan *model.TeamMember)
	// TeamMemberUpdated = make(chan *model.TeamMember)
	// TeamMemberRemoved = make(chan string)
	// TeamCollectionAdded = make(chan *model.TeamCollection)
	// TeamCollectionUpdated = make(chan *model.TeamCollection)
	// TeamCollectionRemoved = make(chan string)
	// TeamCollectionMoved = make(chan *model.TeamCollection)
	// CollectionOrderUpdated = make(chan *dto.CollectionReorderData)
	// TeamRequestAdded = make(chan *model.TeamRequest)
	// TeamRequestUpdated = make(chan *model.TeamRequest)
	// TeamRequestDeleted = make(chan string)
	// RequestOrderUpdated = make(chan *dto.RequestReorderData)
	// RequestMoved = make(chan *model.TeamRequest)
	// TeamInvitationAdded = make(chan *model.TeamInvitation)
	// TeamInvitationRemoved = make(chan string)
	// TeamEnvironmentUpdated = make(chan *model.TeamEnvironment)
	// TeamEnvironmentCreated = make(chan *model.TeamEnvironment)
	// TeamEnvironmentDeleted = make(chan *model.TeamEnvironment)
	// UserSettingsCreated = make(chan *model.UserSetting)
	// UserSettingsUpdated = make(chan *model.UserSetting)
	// UserEnvironmentCreated = make(chan *model.UserEnvironment)
	// UserEnvironmentUpdated = make(chan *model.UserEnvironment)
	// UserEnvironmentDeleted = make(chan *model.UserEnvironment)
	// UserEnvironmentDeleteMany = make(chan int)
	// UserHistoryCreated = make(chan *model.UserHistory)
	// UserHistoryUpdated = make(chan *model.UserHistory)
	// UserHistoryDeleted = make(chan *model.UserHistory)
	// UserHistoryDeletedMany = make(chan *dto.UserHistoryDeletedManyData)
	// UserRequestCreated = make(chan *model.UserRequest)
	// UserRequestUpdated = make(chan *model.UserRequest)
	// UserRequestDeleted = make(chan *model.UserRequest)
	// UserRequestMoved = make(chan *dto.UserRequestReorderData)
	// UserCollectionCreated = make(chan *model.UserCollection)
	// UserCollectionUpdated = make(chan *model.UserCollection)
	// UserCollectionRemoved = make(chan *dto.UserCollectionRemovedData)
	// UserCollectionMoved = make(chan *model.UserCollection)
	// UserCollectionOrderUpdated = make(chan *dto.UserCollectionReorderData)
	// MyShortcodesCreated = make(chan *model.Shortcode)
	// MyShortcodesUpdated = make(chan *model.Shortcode)
	// MyShortcodesRevoked = make(chan *model.Shortcode)

	UserUpdatedSub = pubsub.NewPublisher(time.Second, 10)
	UserDeletedSub = pubsub.NewPublisher(time.Second, 10)
	UserInvitedSub = pubsub.NewPublisher(time.Second, 10)
	TeamMemberAddedSub = pubsub.NewPublisher(time.Second, 10)
	TeamMemberUpdatedSub = pubsub.NewPublisher(time.Second, 10)
	TeamMemberRemovedSub = pubsub.NewPublisher(time.Second, 10)
	TeamCollectionAddedSub = pubsub.NewPublisher(time.Second, 10)
	TeamCollectionUpdatedSub = pubsub.NewPublisher(time.Second, 10)
	TeamCollectionRemovedSub = pubsub.NewPublisher(time.Second, 10)
	TeamCollectionMovedSub = pubsub.NewPublisher(time.Second, 10)
	CollectionOrderUpdatedSub = pubsub.NewPublisher(time.Second, 10)
	TeamRequestAddedSub = pubsub.NewPublisher(time.Second, 10)
	TeamRequestUpdatedSub = pubsub.NewPublisher(time.Second, 10)
	TeamRequestDeletedSub = pubsub.NewPublisher(time.Second, 10)
	RequestOrderUpdatedSub = pubsub.NewPublisher(time.Second, 10)
	RequestMovedSub = pubsub.NewPublisher(time.Second, 10)
	TeamInvitationAddedSub = pubsub.NewPublisher(time.Second, 10)
	TeamInvitationRemovedSub = pubsub.NewPublisher(time.Second, 10)
	TeamEnvironmentUpdatedSub = pubsub.NewPublisher(time.Second, 10)
	TeamEnvironmentCreatedSub = pubsub.NewPublisher(time.Second, 10)
	TeamEnvironmentDeletedSub = pubsub.NewPublisher(time.Second, 10)
	UserSettingsCreatedSub = pubsub.NewPublisher(time.Second, 10)
	UserSettingsUpdatedSub = pubsub.NewPublisher(time.Second, 10)
	UserEnvironmentCreatedSub = pubsub.NewPublisher(time.Second, 10)
	UserEnvironmentUpdatedSub = pubsub.NewPublisher(time.Second, 10)
	UserEnvironmentDeletedSub = pubsub.NewPublisher(time.Second, 10)
	UserEnvironmentDeleteManySub = pubsub.NewPublisher(time.Second, 10)
	UserHistoryCreatedSub = pubsub.NewPublisher(time.Second, 10)
	UserHistoryUpdatedSub = pubsub.NewPublisher(time.Second, 10)
	UserHistoryDeletedSub = pubsub.NewPublisher(time.Second, 10)
	UserHistoryDeletedManySub = pubsub.NewPublisher(time.Second, 10)
	UserRequestCreatedSub = pubsub.NewPublisher(time.Second, 10)
	UserRequestUpdatedSub = pubsub.NewPublisher(time.Second, 10)
	UserRequestDeletedSub = pubsub.NewPublisher(time.Second, 10)
	UserRequestMovedSub = pubsub.NewPublisher(time.Second, 10)
	UserCollectionCreatedSub = pubsub.NewPublisher(time.Second, 10)
	UserCollectionUpdatedSub = pubsub.NewPublisher(time.Second, 10)
	UserCollectionRemovedSub = pubsub.NewPublisher(time.Second, 10)
	UserCollectionMovedSub = pubsub.NewPublisher(time.Second, 10)
	UserCollectionOrderUpdatedSub = pubsub.NewPublisher(time.Second, 10)
	MyShortcodesCreatedSub = pubsub.NewPublisher(time.Second, 10)
	MyShortcodesUpdatedSub = pubsub.NewPublisher(time.Second, 10)
	MyShortcodesRevokedSub = pubsub.NewPublisher(time.Second, 10)
}
