import 'package:flutter/material.dart';
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/designkit.dart' as ds;

class SearchMimetypeDropdown extends StatelessWidget {
  final media.MediaSearchRequest current;
  final void Function(media.MediaSearchRequest r) onChange;
  SearchMimetypeDropdown(this.current, {super.key, required this.onChange});

  static Icon icon(int checksum) {
    if (checksum == mimex.checksumfor(mimex.icomovie)) return Icon(Icons.movie_filter);
    if (checksum == mimex.checksumfor(mimex.icoaudio)) return Icon(Icons.music_note);
    return Icon(Icons.file_open_rounded);
  }

  static String label(int checksum) {
    if (checksum == mimex.checksumfor(mimex.icomovie)) return "Movies";
    if (checksum == mimex.checksumfor(mimex.icoaudio)) return "Music";
    return "Files";
  }

  static List<String> mimetypesFor(int checksum) {
    if (checksum == mimex.checksumfor(mimex.icomovie)) return mimex.of(mimex.icomovie);
    if (checksum == mimex.checksumfor(mimex.icoaudio)) return mimex.of(mimex.icoaudio);
    return const [];
  }

  static void select(media.MediaSearchRequest current, int checksum) {
    current.mimetypes
      ..clear()
      ..addAll(mimetypesFor(checksum));
  }

  static List<PopupMenuEntry<String>> menuItems(
    media.MediaSearchState search,
    Function(media.MediaSearchState) onChange,
  ) {
    return [
      _menuOption(search, mimex.checksumfor(mimex.icoaudio), onChange),
      _menuOption(search, mimex.checksumfor(mimex.icomovie), onChange),
      _menuOption(search, mimex.checksumfor(mimex.icobinary), onChange),
    ];
  }

  static PopupMenuItem<String> _menuOption(
    media.MediaSearchState search,
    int checksum,
    Function(media.MediaSearchState) onChange,
  ) {
    final selected = mimex.checksum(search.next.mimetypes) == checksum;
    return PopupMenuItem<String>(
      enabled: !selected,
      mouseCursor: selected ? SystemMouseCursors.basic : SystemMouseCursors.click,
      onTap: () {
        final current = search.next;
        select(current, checksum);
        onChange(media.MediaSearchState(next: current.clone(), count: search.count));
      },
      child: ds.build((context) {
        final defaults = ds.Defaults.of(context);
        return Row(
          spacing: defaults.spacing,
          children: [
            icon(checksum),
            Text(
              label(checksum),
              style: TextStyle(fontWeight: selected ? FontWeight.bold : FontWeight.normal),
            ),
            const Spacer(),
          ],
        );
      }),
    );
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
        icon: icon(mimetypes),
        onSelected: (v) {
          select(current, v);
          onChange(current);
        },
        itemBuilder: (context) => [
          PopupMenuItem(
            value: mimex.checksumfor(mimex.icoaudio),
            child: Tooltip(message: "Music", child: Icon(Icons.music_note)),
          ),
          PopupMenuItem(
            value: mimex.checksumfor(mimex.icomovie),
            child: Tooltip(
              message: "Movies",
              child: Icon(Icons.movie_filter),
            ),
          ),
          PopupMenuItem(
            value: mimex.checksumfor(mimex.icobinary),
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
