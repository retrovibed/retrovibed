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

  static List<PopupMenuEntry<String>> menuItems(ValueNotifier<media.MediaSearchState> search) {
    return [
      _menuOption(search, mimex.checksumfor(mimex.icoaudio)),
      _menuOption(search, mimex.checksumfor(mimex.icomovie)),
      _menuOption(search, mimex.checksumfor(mimex.icobinary)),
    ];
  }

  static PopupMenuItem<String> _menuOption(ValueNotifier<media.MediaSearchState> search, int checksum) {
    return PopupMenuItem<String>(
      child: ValueListenableBuilder<media.MediaSearchState>(
        valueListenable: search,
        builder: (context, state, _) {
          final selected = mimex.checksum(state.next.mimetypes) == checksum;
          return ListTile(
            leading: icon(checksum),
            title: Text(label(checksum)),
            selected: selected,
            enabled: !selected,
            onTap: () {
              final current = state.next;
              select(current, checksum);
              search.value = media.MediaSearchState(next: current.clone(), count: search.value.count);
            },
          );
        },
      ),
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
        itemBuilder:
            (context) => [
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
