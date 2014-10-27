package openflow

import (
	"errors"

	"github.com/kandoo/beehive-netctrl/openflow/of10"
	"github.com/kandoo/beehive-netctrl/openflow/of12"
)

func (of *of10Driver) handleFeaturesReply(rep of10.FeaturesReply,
	c *ofConn) error {
	return lateFeaturesReplyError()
}

func (of *of12Driver) handleFeaturesReply(rep of12.FeaturesReply,
	c *ofConn) error {
	return lateFeaturesReplyError()
}

func lateFeaturesReplyError() error {
	return errors.New("Cannot receive a features reply after handshake.")
}
