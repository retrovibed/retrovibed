import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/discovery.dart' as disc;
import 'search.dart';

class Home extends StatefulWidget {
  final media.FnMediaSearch apisearch;
  final media.FnUploadRequest apiupload;
  final TextEditingController? controller;
  final FocusNode? focus;
  final String highlighted;
  final ValueNotifier<media.MediaSearchState> search;

  const Home({
    super.key,
    this.apisearch = media.media.search,
    this.apiupload = media.media.upload,
    this.controller,
    this.focus,
    required this.highlighted,
    required this.search,
  });

  @override
  State<StatefulWidget> createState() => _HomeState();
}

enum _Mode { library, discovery, remote }

class _HomeState extends State<Home> {
  Widget _downloading = ds.Empty;
  final ValueNotifier<_Mode> _mode = ValueNotifier(_Mode.library);

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _switchToMode(_Mode m) {
    _mode.value = m;
    widget.focus?.requestFocus();
  }

  @override
  void dispose() {
    _mode.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder<_Mode>(
      valueListenable: _mode,
      builder: (context, mode, _) => switch (mode) {
        _Mode.library => LibrarySearch(
          apisearch: widget.apisearch,
          apiupload: widget.apiupload,
          controller: widget.controller,
          focus: widget.focus,
          highlighted: widget.highlighted,
          search: widget.search,
          discovering: false,
          onToggleMode: () => _switchToMode(_Mode.library),
          downloading: _downloading,
          onDownloadingChanged: (w) => setState(() => _downloading = w),
        ),
        _Mode.discovery => disc.DiscoverySearch(
          apiupload: widget.apiupload,
          controller: widget.controller,
          focus: widget.focus,
          search: widget.search,
          discovering: true,
          onToggleMode: () => _switchToMode(_Mode.discovery),
          downloading: _downloading,
          onDownloadingChanged: (w) => setState(() => _downloading = w),
        ),
        _Mode.remote => remote.Display(),
      },
    );
  }
}
