import 'dart:async';

import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/filesystem/api.dart' as api;
import 'package:retrovibed/library/list.display.dart' as lib;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/uuidx.dart' as uuidx;

// moves an entry into a directory picked from every directory in the library. the picker
// is the library list fixed to the directory mimetype, so it searches rather than walks.
class ButtonFileMove extends StatelessWidget {
  final media.Media current;
  // null signals the entry left the listed directory.
  final void Function(media.Media? upd)? onChange;
  final media.FnMediaSearch apisearch;
  final api.FnFilesystemMove apimove;

  const ButtonFileMove({
    super.key,
    required this.current,
    this.onChange,
    this.apisearch = media.media.search,
    this.apimove = api.filesystem.move,
  });

  @override
  Widget build(BuildContext context) {
    return ds.LoadingIconButton(
      icon: const Icon(Icons.drive_file_move_outline),
      tooltip: "move to another folder",
      help: ds.Hint(const Text("move to another folder")),
      onPressed: () {
        // captured while the caller's context is still valid; the picker outlives it.
        final options = [authn.request(authn.AuthzCache.meta(context))];
        return ds.modals
            .asyncfn<media.Media?>(
              context,
              (completion) => ConstrainedBox(
                constraints: const BoxConstraints(maxWidth: 512.0),
                child: _FileMovePicker(
                  current,
                  completion: completion,
                  options: options,
                  apisearch: apisearch,
                  apimove: apimove,
                ),
              ),
            )
            .then((moved) {
              if (moved != null) onChange?.call(null);
            });
      },
    );
  }
}

class _FileMovePicker extends StatefulWidget {
  final media.Media current;
  final Completer<media.Media?> completion;
  final List<httpx.Option> options;
  final media.FnMediaSearch apisearch;
  final api.FnFilesystemMove apimove;

  const _FileMovePicker(
    this.current, {
    required this.completion,
    required this.options,
    required this.apisearch,
    required this.apimove,
  });

  @override
  State<StatefulWidget> createState() => _FileMovePickerState();
}

class _FileMovePickerState extends State<_FileMovePicker> with ds.LoadingState {
  Future<void> move(String directory) {
    return httpx
        .withRetry(
          () => widget.apimove(
            widget.current.id,
            api.FilesystemMoveRequest(directoryId: directory),
            options: widget.options,
          ),
        )
        .then((v) => widget.completion.complete(v.media))
        .catchError((cause) {
          setState(() {
            this.cause = ds.Error.unauthorized(cause, onTap: reseterr);
          });
        }, test: httpx.ErrorsTest.unauthorized)
        .catchError((cause) {
          setState(() {
            this.cause = ds.Error.unknown(cause, onTap: reseterr);
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    return ds.Card(
      leading: [
        ds.Heading(
          Text("move ${widget.current.description}"),
          trailing: [
            IconButton(icon: const Icon(Icons.close), onPressed: () => widget.completion.complete(null)),
          ],
        ),
      ],
      ds.Loading(
        cause: cause,
        lib.AvailableListDisplay(
          search: widget.apisearch,
          mimetypes: [mimex.directory],
          leading: [
            media.RowDisplay(
              media: media.Media(id: uuidx.min(), description: "/", mimetype: mimex.directory),
              leading: [Icon(mimex.icofolder)],
              onTap: () => move(uuidx.min()),
            ),
          ],
          // a directory cannot hold itself.
          row: (v) => v.id == widget.current.id
              ? ds.Empty
              : media.RowDisplay(media: v, leading: [Icon(mimex.icofolder)], onTap: () => move(v.id)),
        ),
      ),
    );
  }
}
