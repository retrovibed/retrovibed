import 'dart:io';
import 'package:retrovibed/media.dart';
import 'package:flutter/material.dart';
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/authn.dart' as authn;
import 'package:downloadsfolder/downloadsfolder.dart';
import 'package:path_provider/path_provider.dart';
import 'package:retrovibed/media/content.addressable.storage.pb.dart' as cas;
import './api.dart' as api;

Future<void> Function()? PlayAction(
  BuildContext context,
  Media current,
  MediaSearchResponse s,
) {
  switch (mimex.icon(current.mimetype)) {
    case mimex.icomovie:
    case mimex.icoaudio:
      final playlist = Playlist.of(context);
      return playlist == null
          ? null
          : () {
            return Future.sync(
              () => playlist.setPlaylist(
                s.next,
                range(s.next, current, options: () => [authn.AuthzCache.bearer(context)]),
              ),
            );
          };
    default:
      return null;
  }
}

Future<void> Function() DownloadAction(BuildContext context, Media current) {
  return () {
    if (Platform.isIOS || Platform.isAndroid) {
      return getTemporaryDirectory().then((downloads) async {
        final fname = current.description.trim().replaceAll(" ", ".");
        final dst = File('${downloads.path}/$fname');
        final sink = dst.openWrite();
        return api.media
            .download(current.id, options: [authn.AuthzCache.bearer(context)])
            .then((resp) => resp.stream.pipe(sink))
            .whenComplete(() {
              return sink.close().then((v) {
                return copyFileIntoDownloadFolder(
                  dst.path,
                  fname,
                ).then((b) => b ?? false ? Future.value() : Future.error("failed to copy file $fname"));
              });
            });
      });
    }

    return getDownloadsDirectory().then((downloads) {
      final sink =
          File(
            '${downloads!.path}/${current.description.trim().replaceAll(" ", ".")}',
          ).openWrite();
      return api.media
          .download(current.id, options: [authn.AuthzCache.bearer(context)])
          .then((resp) => resp.stream.pipe(sink))
          .whenComplete(() => sink.close());
    });
  };
}

Future<Media> Function() ArchiveAction(
  BuildContext context,
  Media current, {
  Media Function(Media)? then,
}) {
  return () {
    return api.media
        .update(
          current.id,
          current..archiveId = uuidx.max(),
          options: [authn.AuthzCache.bearer(context)],
        )
        .then((v) => v.media)
        .then(then ?? (v) => v);
  };
}

Future<Media> Function() ArchiveCancelAction(
  BuildContext context,
  Media current, {
  Media Function(Media)? then,
}) {
  return () {
    return api.media
        .update(
          current.id,
          current..archiveId = uuidx.min(),
          options: [authn.AuthzCache.bearer(context)],
        )
        .then((v) => v.media)
        .then(then ?? (v) => v);
  };
}

Future<Media> Function() ArchivePurgeAction(
  BuildContext context,
  Media current, {
  Media Function(cas.Media)? then,
}) {
  return () {
    return api.media
        .unarchive(
          current.archiveId,
          options: [authn.DeeppoolAuthzCache.bearer(context)],
        )
        .then((v) => v.media)
        .then(
          then ?? (v) => current..archiveId = uuidx.min(),
        );
  };
}

class ButtonPlay extends StatelessWidget {
  final Media current;
  final MediaSearchResponse playlist;
  const ButtonPlay({super.key, required this.current, required this.playlist});

  @override
  Widget build(BuildContext context) {
    switch (mimex.icon(current.mimetype)) {
      case mimex.icomovie:
      case mimex.icoaudio:
        return IconButton(
          icon: Icon(Icons.play_circle_outline_rounded),
          onPressed: PlayAction(context, current, playlist),
        );
      default:
        return Container();
    }
  }
}
