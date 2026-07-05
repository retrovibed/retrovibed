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
        return SingleChildScrollView(
          child: ConstrainedBox(
            constraints: BoxConstraints(minHeight: constraints.maxHeight),
            child: Container(
              padding: defaults.padding,
              child: AvailableListDisplay(
                search: widget.availableSearch,
                controller: controller,
                events: refresh,
                trailing: [
                  DownloadingListDisplay(
                    search: widget.downloadingSearch,
                    watch: widget.downloadWatch,
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
