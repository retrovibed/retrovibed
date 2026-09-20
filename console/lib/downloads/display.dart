import 'dart:async';
import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'downloading.list.dart';
import 'available.list.dart';

class Display extends StatefulWidget {
  final media.FnDownloadSearch downloadingSearch;
  final media.FnDownloadSearch apiavailablesearch;
  final media.FnDownloadWatch downloadWatch;
  final List<Widget> leading;
  const Display({
    super.key,
    this.downloadingSearch = media.discovered.downloading,
    this.apiavailablesearch = media.discovered.available,
    this.downloadWatch = media.discovered.watch,
    this.leading = const [],
  });

  @override
  State<Display> createState() => _DisplayState();
}

class _DisplayState extends State<Display> {
  final TextEditingController controller = TextEditingController();
  final StreamController<media.Download> refresh = StreamController<media.Download>.broadcast();

  @override
  void dispose() {
    controller.dispose();
    refresh.close();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return ds.Container(
      AvailableListDisplay(
        refresh,
        search: widget.apiavailablesearch,
        controller: controller,
        leading: widget.leading,
        trailing: [
          DownloadingListDisplay(
            refresh,
            search: widget.downloadingSearch,
            watch: widget.downloadWatch,
          ),
        ],
      ),
    );
  }
}
