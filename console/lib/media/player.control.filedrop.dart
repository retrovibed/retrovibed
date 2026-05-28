import 'package:flutter/material.dart';
import 'package:media_kit/media_kit.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/mimex.dart' as mimex;
import './playlist.dart' as internal;

class PlayerControlFiledrop extends StatelessWidget {
  final Player player;
  const PlayerControlFiledrop(this.player, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    const mimetypes = <String>[
      // ...mimex.folders,
      ...mimex.audios,
      ...mimex.videos,
      "application/x-iso9660-image", // sometimes bluray/dvd are iso images.
    ];
    final playfile = (ds.FilesEvent evt, {ValueNotifier<int>? progress}) {
      print("play file checkpoint 0 ${evt.files.length}");
      if (evt.files.isEmpty) return Future.value(ds.NullWidget);
      final file = evt.files.firstWhere((v) {
        print("play file checkpoint 1 ${v.mimeType}");
        return mimetypes.any((x) => x == v.mimeType);
      }, orElse: () => evt.files.first);

      print("play file checkpoint 2 ${file.name}");
      return internal.Playlist.file(context, "file://${file.path}").then((v) => ds.NullWidget);
    };

    return ds.FileDropWell.icon(
      playfile,
      mimetypes: mimetypes,
      icon: Icons.video_collection_rounded,
      tooltip: "play a local media file",
      help: ds.Hint(const Text("play a local media file")),
    );
  }
}
