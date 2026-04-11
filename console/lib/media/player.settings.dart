import 'package:flutter/material.dart';
import 'package:media_kit/media_kit.dart'; // Provides [Player], [Media], [Playlist] etc.
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/designkit.dart' as ds;

class PlayerSettings extends StatefulWidget {
  final Player current;
  final EdgeInsets? padding;
  final EdgeInsets? margin;
  final BoxDecoration? decoration;
  final BoxConstraints? constraints;
  final Alignment? alignment;
  final Color? background;

  const PlayerSettings({
    super.key,
    required this.current,
    this.padding,
    this.margin,
    this.decoration,
    this.constraints,
    this.alignment,
    this.background,
  });

  @override
  State<PlayerSettings> createState() => _PlayerSettingsState();
}

class _PlayerSettingsState extends State<PlayerSettings> {
  Track _track = Track();

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  @override
  void initState() {
    super.initState();
    setState(() {
      _track = widget.current.state.track;
    });
  }

  static String? audio(AudioTrack v) {
    return v.title ?? v.language;
  }

  static String? video(VideoTrack v) {
    return v.title ?? v.language;
  }

  static String? subtitle(SubtitleTrack v) {
    return v.title ?? v.language;
  }

  @override
  Widget build(BuildContext context) {
    final _audiolist =
        widget.current.state.tracks.audio.map((v) {
          return DropdownMenuItem(
            child: Text(_PlayerSettingsState.audio(v) ?? v.id),
            value: v,
          );
        }).toList();
    final _videolist =
        widget.current.state.tracks.video.map((v) {
          return DropdownMenuItem(
            child: Text(_PlayerSettingsState.video(v) ?? v.id),
            value: v,
          );
        }).toList();
    final _subtitlelist =
        widget.current.state.tracks.subtitle.map((v) {
          return DropdownMenuItem(
            child: Text(_PlayerSettingsState.subtitle(v) ?? v.id),
            value: v,
          );
        }).toList();

    return ds.Container(
      Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          forms.Field(
            label: Text("Audio"),
            input: DropdownButton(
              alignment: Alignment.topLeft,
              isExpanded: true,
              value: _track.audio,
              items: _audiolist,
              onChanged: (v) {
                setState(() {
                  _track = _track.copyWith(audio: v);
                });
                widget.current.setAudioTrack(v ?? AudioTrack.auto());
              },
            ),
          ),
          forms.Field(
            label: Text("Video"),
            input: DropdownButton(
              alignment: Alignment.topLeft,
              isExpanded: true,
              value: _track.video,
              items: _videolist,
              onChanged: (v) {
                setState(() {
                  _track = _track.copyWith(video: v);
                });
                widget.current.setVideoTrack(v ?? VideoTrack.auto());
              },
            ),
          ),
          forms.Field(
            label: Text("Subtitle"),
            input: DropdownButton(
              alignment: Alignment.topLeft,
              isExpanded: true,
              value: _track.subtitle,
              items: _subtitlelist,
              onChanged: (v) {
                setState(() {
                  _track = _track.copyWith(subtitle: v);
                });
                widget.current.setSubtitleTrack(v ?? SubtitleTrack.auto());
              },
            ),
          ),
        ],
      ),
      padding: widget.padding,
      margin: widget.margin,
      decoration: widget.decoration,
      constraints: widget.constraints,
      alignment: widget.alignment,
      background: widget.background,
    );
  }
}
