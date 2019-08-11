package main
import (
  "errors"
)
// ErrNoAvatar occurs when Avatar instance
// cannot return URL
var ErrNoAvatarURL = errors.New("chat: Unable to get avatar url ")
// Avatar shows profile image
type Avatar interface {
  // GetAvatarURL returns client's avatar url
  // return ErrNoAvater when unable to get url
  GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}
var UseAuthAvatar AuthAvatar
func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
  if url, ok := c.userData["avatar_url"]; ok {
    if urlStr, ok := url.(string); ok {
      return urlStr, nil
    }
  }
  return "", ErrNoAvatarURL
}
