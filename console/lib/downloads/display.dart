import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'downloading.list.dart';
import 'available.list.dart';

class Display extends StatefulWidget {
  final media.FnDownloadSearch downloadingSearch;
  final media.FnDownloadSearch availableSearch;
  final media.FnDownloadWatch downloadWatch;
  const Display({
    super.key,
    this.downloadingSearch = media.discovered.downloading,
    this.availableSearch = media.discovered.available,
    this.downloadWatch = media.discovered.watch,
  });
  static const hints = [
    ds.Hint(
      label: const Text("Progress"),
      description: const Text(
        "active downloads auto-refresh at the top of the view",
      ),
    ),
    ds.Hint(
      label: const Text("Search"),
      description: const Text("find available content to download"),
    ),
    ds.Hint(
      label: const Text("Upload"),
      description: const Text("drag and drop files to discover new content"),
    ),
    ds.Hint(
      label: const Text("Magnet"),
      description: const Text("import magnet links via the magnet button"),
    ),
    ds.Hint(
      label: const Text("Double-click"),
      description: const Text(
        "expand a download to see size, path, and status details",
      ),
    ),
  ];

  @override
  State<Display> createState() => _DisplayState();
}

class _DisplayState extends State<Display> {
  final TextEditingController controller = TextEditingController();
  final ValueNotifier<int> refresh = ValueNotifier(0);

  @override
  void dispose() {
    controller.dispose();
    refresh.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final defaults = ds.Defaults.of(context);
        final compact = defaults.isCompact;
        return SingleChildScrollView(
          child: ConstrainedBox(
            constraints: BoxConstraints(minHeight: constraints.maxHeight),
            child: Container(
              padding: defaults.padding,
              child: Column(
                mainAxisSize: MainAxisSize.max,
                mainAxisAlignment: MainAxisAlignment.start,
                verticalDirection: compact ? VerticalDirection.up : VerticalDirection.down,
                spacing: 0.0,
                children: [
                  DownloadingListDisplay(
                    search: widget.downloadingSearch,
                    watch: widget.downloadWatch,
                    events: refresh,
                  ),
                  AvailableListDisplay(
                    search: widget.availableSearch,
                    controller: controller,
                    events: refresh,
                  ),
                ],
              ),
            ),
          ),
        );
      },
    );
  }
}
