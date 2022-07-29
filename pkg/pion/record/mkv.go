package record

import (
	"bytes"
	"os"
	"time"

	"github.com/at-wat/ebml-go/mkvcore"
	"github.com/at-wat/ebml-go/webm"
	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
	"github.com/pion/webrtc/v3/pkg/media/samplebuilder"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type MkvSaver struct {
	audioWriter, videoWriter       webm.BlockWriteCloser
	audioBuilder, videoBuilder     *samplebuilder.SampleBuilder
	audioTimestamp, videoTimestamp time.Duration
	onlyAudioWriter                *oggwriter.OggWriter
	savePath                       string
}

func NewMkvSaver(savePath string) *MkvSaver {
	return &MkvSaver{
		audioBuilder: samplebuilder.New(32, &codecs.OpusPacket{}, 48000),
		videoBuilder: samplebuilder.New(1024, &codecs.H264Packet{}, 90000),
		savePath:     savePath,
	}
}

func (s *MkvSaver) InitWriter(width, height int) {
	if width == 0 || height == 0 {
		var err error
		s.onlyAudioWriter, err = oggwriter.New(s.savePath+".ogg", 48000, 2)
		if err != nil {
			logrus.Errorf("fail to open %s", s.savePath+".mkv")
		}
		return
	}
	if s.onlyAudioWriter != nil {
		s.onlyAudioWriter.Close()
	}
	w, err := os.OpenFile(s.savePath+".mkv", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		logrus.Error(err)
		return
	}
	var desc []mkvcore.TrackDescription
	desc = append(desc,
		mkvcore.TrackDescription{
			TrackNumber: uint64(1),
			TrackEntry: webm.TrackEntry{
				Name:        "Audio",
				TrackNumber: 1,
				CodecID:     "A_OPUS",
				TrackType:   2,
				Audio: &webm.Audio{
					SamplingFrequency: 48000.0,
					Channels:          2,
				},
			},
		},
	)

	desc = append(desc,
		mkvcore.TrackDescription{
			TrackNumber: uint64(2),
			TrackEntry: webm.TrackEntry{
				Name:        "Video",
				TrackNumber: 2,
				CodecID:     "V_MPEG4/ISO/AVC",
				TrackType:   1,
				Video: &webm.Video{
					PixelWidth:  uint64(width),
					PixelHeight: uint64(height),
				},
			},
		},
	)
	header := webm.DefaultEBMLHeader
	header.DocType = "matroska"
	ws, err := mkvcore.NewSimpleBlockWriter(
		w, desc,
		mkvcore.WithEBMLHeader(header),
		mkvcore.WithSegmentInfo(webm.DefaultSegmentInfo),
		mkvcore.WithBlockInterceptor(webm.DefaultBlockInterceptor),
	)

	s.audioWriter = ws[0]
	s.videoWriter = ws[1]
}

func (s *MkvSaver) Push264(rtpPacket *rtp.Packet) {
	s.videoBuilder.Push(rtpPacket)

	for {
		sample := s.videoBuilder.Pop()
		if sample == nil {
			return
		}
		naluType := sample.Data[4] & 0x1F
		videoKeyframe := (naluType == 7) || (naluType == 8)
		if videoKeyframe {
			if (s.videoWriter == nil || s.audioWriter == nil) && naluType == 7 {
				p := bytes.SplitN(sample.Data[4:], []byte{0x00, 0x00, 0x00, 0x01}, 2)
				if width, height, fps, ok := H264_decode_sps(p[0], uint(len(p[0]))); ok {
					logrus.Infof("width:%d, height:%d, fps:%d", width, height, fps)
					s.InitWriter(width, height)
				}
			}
		}
		if s.videoWriter != nil {
			s.videoTimestamp += sample.Duration
			if _, err := s.videoWriter.Write(videoKeyframe, int64(s.videoTimestamp/time.Millisecond), sample.Data); err != nil {
				return

			}
		}
	}
}

func (s *MkvSaver) PushOpus(rtpPacket *rtp.Packet) {
	if s.onlyAudioWriter != nil {
		s.onlyAudioWriter.WriteRTP(rtpPacket)
		return
	}
	s.audioBuilder.Push(rtpPacket)

	for {

		sample := s.audioBuilder.Pop()
		if sample == nil {
			return
		}
		if s.audioWriter != nil {
			s.audioTimestamp += sample.Duration
			if _, err := s.audioWriter.Write(true, int64(s.videoTimestamp/time.Millisecond), sample.Data); err != nil {
				return
			}
		}
	}
}

func (s *MkvSaver) Close() error {
	if s.onlyAudioWriter != nil {
		if err := s.onlyAudioWriter.Close(); err != nil {
			return errors.Wrapf(err, "关闭音频失败")
		}
	}
	if s.audioWriter != nil {
		if err := s.audioWriter.Close(); err != nil {
			return errors.Wrapf(err, "关闭音频失败")
		}
	}
	if s.videoWriter != nil {
		if err := s.videoWriter.Close(); err != nil {
			return errors.Wrapf(err, "关闭视频失败")
		}
	}
	return nil
}
