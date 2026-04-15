import 'package:flutter/material.dart';
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/designkit.dart' as ds;

class SearchMimetypeDropdown extends StatelessWidget {
  final media.MediaSearchRequest current;
  final void Function(media.MediaSearchRequest r) onChange;
  SearchMimetypeDropdown(this.current, {super.key, required this.onChange});

  Icon _icon(int checksum) {
    if (checksum == mimex.checksumfor(mimex.movie)) return Icon(Icons.movie_filter);
    if (checksum == mimex.checksumfor(mimex.audio)) return Icon(Icons.music_note);
    return Icon(Icons.file_open_rounded);
  }

  @override
  Widget build(BuildContext context) {
    final mimetypes = mimex.checksum(current.mimetypes);

    return ds.Help(
      PopupMenuButton<int>(
        tooltip: "filter by mimetype",
        position: PopupMenuPosition.under,
        color: Theme.of(context).colorScheme.surface,
        surfaceTintColor: Theme.of(context).colorScheme.surface,
        icon: _icon(mimetypes),
        onSelected: (v) {
          if (v == mimex.checksumfor(mimex.movie)) {
            current.mimetypes.clear();
            current.mimetypes.addAll(mimex.of(mimex.movie));
            onChange(current);
            return;
          }

          if (v == mimex.checksumfor(mimex.audio)) {
            current.mimetypes.clear();
            current.mimetypes.addAll(mimex.of(mimex.audio));
            onChange(current);
            return;
          }

          current.mimetypes.clear();
          onChange(current);
        },
        itemBuilder:
            (context) => [
              PopupMenuItem(
                value: mimex.checksumfor(mimex.audio),
                child: Tooltip(message: "Music", child: Icon(Icons.music_note)),
              ),
              PopupMenuItem(
                value: mimex.checksumfor(mimex.movie),
                child: Tooltip(
                  message: "Movies",
                  child: Icon(Icons.movie_filter),
                ),
              ),
              PopupMenuItem(
                value: mimex.checksumfor(mimex.binary),
                child: Tooltip(
                  message: "Files",
                  child: Icon(Icons.file_open_rounded),
                ),
              ),
            ],
      ),
      ds.Hint(const Text("dropdown to filter by movies, audio, documents, or images")),
    );
  }
}
