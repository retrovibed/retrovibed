import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/downloads.dart' as downloads;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'dropdown.nav.menu.dart';
import 'search.mimetype.dropdown.dart';

// Shared mimetype-filter + mode-switch + upload dropdown used by both the
// library and discovery search bars.
class SearchUploadDropdown extends StatelessWidget {
  final ValueNotifier<media.MediaSearchState> search;
  final ValueNotifier<media.SearchMode> mode;
  final void Function(media.SearchMode) onModeChanged;
  final media.FnUploadRequest apiupload;
  final void Function(Widget) onDownloadingChanged;

  const SearchUploadDropdown({
    super.key,
    required this.search,
    required this.mode,
    required this.onModeChanged,
    required this.apiupload,
    required this.onDownloadingChanged,
  });

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder<media.MediaSearchState>(
      valueListenable: search,
      builder: (context, state, _) => DropdownNavMenu.options(
        icon: SearchMimetypeDropdown.icon(mimex.checksum(state.next.mimetypes)),
        help: ds.Hint(
          const Text(
            "filter by mimetype, upload files, torrents, magnet links, or switch between library, files, discover, and downloads mode",
          ),
        ),
        search: search,
        mode: mode,
        onModeChanged: onModeChanged,
        options: [
          PopupMenuItem<String>(
            enabled: false,
            child: ValueListenableBuilder<media.MediaSearchState>(
              valueListenable: search,
              builder: (context, s, _) => mimex.CategoryOptionsLabel(s.next.mimetypes),
            ),
          ),
          media.MenuItemUploadFiles(
            context,
            search,
            apiupload: apiupload,
          ),
          downloads.MenuItemDownloadTorrent(context, (d) {
            onDownloadingChanged(
              downloads.DownloadQueue(
                d,
                onQueueComplete: () => onDownloadingChanged(ds.Empty),
              ),
            );
            print("downloading torrents ${d}");
          }),
          downloads.MenuItemDownloadMagnet(context, (d) {
            onDownloadingChanged(
              downloads.DownloadQueue(
                d,
                onQueueComplete: () => onDownloadingChanged(ds.Empty),
              ),
            );
            print("downloading magnets ${d}");
          }),
        ],
      ),
    );
  }
}
