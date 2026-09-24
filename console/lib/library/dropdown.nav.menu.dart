import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'search.mimetype.dropdown.dart';

class DropdownNavMenu extends StatelessWidget {
  static const Widget defaultHelp = ds.Hint(
    Text("filter by mimetype, or switch between library, files, discover, and downloads mode"),
  );

  final Widget icon;
  final List<PopupMenuEntry<String>> items;
  final Widget help;

  const DropdownNavMenu({
    super.key,
    required this.icon,
    required this.items,
    this.help = defaultHelp,
  });

  // mimetype filters + mode switches, followed by caller specific options after a divider.
  static Widget options({
    Key? key,
    required Widget icon,
    Widget help = defaultHelp,
    required ValueNotifier<media.MediaSearchState> search,
    required ValueNotifier<media.SearchMode> mode,
    required void Function(media.SearchMode) onModeChanged,
    List<PopupMenuEntry<String>> options = const [],
  }) {
    return ds.build(
      (context) => DropdownNavMenu(
        key: key,
        icon: icon,
        help: help,
        items: [
          ...SearchMimetypeDropdown.menuItems(mode.value, search.value, (upd) {
            search.value = upd;
            onModeChanged(media.SearchMode.library);
          }),
          media.SearchModeToggle(
            mode: media.SearchMode.filesystem,
            current: mode,
            icon: mimex.icofolder,
            label: "Files",
            onSelect: onModeChanged,
          ),
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
          if (authn.developer(context).alpha)
            media.SearchModeToggle(
              mode: media.SearchMode.social,
              current: mode,
              icon: Icons.share,
              label: "Social",
              onSelect: onModeChanged,
            ),
          if (options.isNotEmpty) const PopupMenuDivider(),
          ...options,
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return ds.Help(
      PopupMenuButton<String>(
        position: PopupMenuPosition.under,
        color: Theme.of(context).colorScheme.surface,
        surfaceTintColor: Theme.of(context).colorScheme.surface,
        icon: icon,
        itemBuilder: (context) => items,
      ),
      help,
    );
  }
}
