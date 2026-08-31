import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/downloads.dart' as downloads;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'dropdown.upload.dart';
import 'search.mimetype.dropdown.dart';

// Shared mimetype-filter + mode-switch + upload dropdown used by both the
// library and discovery search bars; they differ only in help text.
class SearchUploadDropdown extends StatelessWidget {
  final ValueNotifier<media.MediaSearchState> search;
  final ValueNotifier<media.SearchMode> mode;
  final void Function(media.SearchMode) onModeChanged;
  final media.FnUploadRequest apiupload;
  final void Function(Widget) onDownloadingChanged;
  final Widget help;

  const SearchUploadDropdown({
    super.key,
    required this.search,
    required this.mode,
    required this.onModeChanged,
    required this.apiupload,
    required this.onDownloadingChanged,
    required this.help,
  });

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder<media.MediaSearchState>(
      valueListenable: search,
      builder: (context, state, _) => DropdownUpload(
        icon: SearchMimetypeDropdown.icon(mimex.checksum(state.next.mimetypes)),
        help: help,
        items: [
          ...SearchMimetypeDropdown.menuItems(search),
          media.SearchModeToggle(
            mode: media.SearchMode.discovery,
            current: mode,
            icon: Icons.travel_explore,
            label: "Discover",
            onSelect: onModeChanged,
          ),
          media.SearchModeToggle(
            mode: media.SearchMode.downloads,
            current: mode,
            icon: Icons.download,
            label: "Downloads",
            onSelect: onModeChanged,
          ),
          media.SearchModeToggle(
            mode: media.SearchMode.social,
            current: mode,
            icon: Icons.share,
            label: "Social",
            onSelect: onModeChanged,
          ),
          const PopupMenuDivider(),
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
              media.DownloadQueue(
                d,
                onQueueComplete: () => onDownloadingChanged(ds.Empty),
              ),
            );
            print("downloading torrents ${d}");
          }),
          downloads.MenuItemDownloadMagnet(context, (d) {
            onDownloadingChanged(
              media.DownloadQueue(
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
