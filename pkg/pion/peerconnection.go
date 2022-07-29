package pion

import (
	"github.com/pion/interceptor"
	"github.com/pion/sdp/v3"
	"github.com/pion/webrtc/v3"
)

func NewPeerConnection(configuration webrtc.Configuration, enableTWCC, enableAudioLevel bool) (*webrtc.PeerConnection, error) {
	engine := &webrtc.MediaEngine{}
	if err := engine.RegisterDefaultCodecs(); err != nil {
		return nil, err
	}

	for _, extension := range []string{sdp.SDESMidURI, sdp.SDESRTPStreamIDURI, sdp.TransportCCURI} {
		if extension == sdp.TransportCCURI && !enableTWCC {
			continue
		}
		if err := engine.RegisterHeaderExtension(webrtc.RTPHeaderExtensionCapability{URI: extension}, webrtc.RTPCodecTypeVideo); err != nil {
			return nil, err
		}
	}

	for _, extension := range []string{sdp.SDESMidURI, sdp.SDESRTPStreamIDURI, sdp.AudioLevelURI} {
		if extension == sdp.AudioLevelURI && !enableAudioLevel {
			continue
		}
		if err := engine.RegisterHeaderExtension(webrtc.RTPHeaderExtensionCapability{URI: extension}, webrtc.RTPCodecTypeAudio); err != nil {
			return nil, err
		}
	}
	interceptor := &interceptor.Registry{}
	if err := webrtc.RegisterDefaultInterceptors(engine, interceptor); err != nil {
		return nil, err
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(engine), webrtc.WithInterceptorRegistry(interceptor))
	return api.NewPeerConnection(configuration)
}
