import 'dart:async';
import 'package:flutter/material.dart';
import 'package:media_kit/media_kit.dart';
import 'package:retrovibed/designkit.dart' as ds;
import './playlist.dart' as internal;

class _StreamUrlPrompt extends StatelessWidget {
  final TextEditingController controller;
  final void Function(String) onSubmitted;

  _StreamUrlPrompt({required this.onSubmitted, TextEditingController? controller})
    : controller = controller ?? TextEditingController();

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return Container(
      padding: defaults.padding,
      color: defaults.opaque,
      child: TextField(
        autofocus: true,
        controller: controller,
        decoration: const InputDecoration(
          labelText: 'Enter a stream URL',
          border: OutlineInputBorder(),
        ),
        onSubmitted: onSubmitted,
      ),
    );
  }
}

Future<void> playStream(BuildContext context) {
  return ds.modals.asyncfn<String?>(context, (completion) {
    return _StreamUrlPrompt(onSubmitted: completion.complete);
  }).then((url) {
    final trimmed = url?.trim() ?? "";
    if (trimmed.isEmpty) return Future.value();
    return internal.Playlist.file(context, trimmed);
  });
}

class PlayerControlStream extends StatelessWidget {
  final Player player;
  const PlayerControlStream(this.player, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ds.LoadingIconButton(
      onPressed: () => playStream(context),
      icon: const Icon(Icons.stream),
      tooltip: "play an http(s) stream",
      help: ds.Hint(const Text("play an http(s) stream")),
    );
  }
}
