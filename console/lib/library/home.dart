import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/discovery.dart' as disc;
import 'package:retrovibed/downloads.dart' as downloads;
import 'package:retrovibed/community/social.home.dart' as social;
import 'dropdown.upload.dart';
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

class _HomeState extends State<Home> {
  Widget _downloading = ds.Empty;
  final ValueNotifier<media.SearchMode> _mode = ValueNotifier(media.SearchMode.library);

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _switchToMode(media.SearchMode m) {
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
    return ValueListenableBuilder<media.SearchMode>(
      valueListenable: _mode,
      builder: (context, mode, _) => switch (mode) {
        media.SearchMode.library => Search(
          apisearch: widget.apisearch,
          apiupload: widget.apiupload,
          controller: widget.controller,
          focus: widget.focus,
          highlighted: widget.highlighted,
          search: widget.search,
          mode: _mode,
          onModeChanged: _switchToMode,
          downloading: _downloading,
          onDownloadingChanged: (w) => setState(() => _downloading = w),
        ),
        media.SearchMode.discovery => disc.Search(
          apiupload: widget.apiupload,
          controller: widget.controller,
          focus: widget.focus,
          search: widget.search,
          mode: _mode,
          onModeChanged: _switchToMode,
          downloading: _downloading,
          onDownloadingChanged: (w) => setState(() => _downloading = w),
        ),
        media.SearchMode.downloads => downloads.AutoHelp(
          downloads.MeteredWarning(
            downloads.Display(
              leading: [
                ds.CompactingMenu.pinned(
                  DropdownUpload(
                    icon: const Icon(Icons.download),
                    help: ds.Hint(const Text("switch to library, discover, or social mode")),
                    items: [
                      media.SearchModeToggle(
                        mode: media.SearchMode.library,
                        current: _mode,
                        icon: Icons.video_library,
                        label: "Library",
                        onSelect: _switchToMode,
                      ),
                      media.SearchModeToggle(
                        mode: media.SearchMode.discovery,
                        current: _mode,
                        icon: Icons.travel_explore,
                        label: "Discover",
                        onSelect: _switchToMode,
                      ),
                      media.SearchModeToggle(
                        mode: media.SearchMode.social,
                        current: _mode,
                        icon: Icons.share,
                        label: "Social",
                        onSelect: _switchToMode,
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
        media.SearchMode.social => social.SocialHome(
          mode: _mode,
          onModeChanged: _switchToMode,
        ),
      },
    );
  }
}
