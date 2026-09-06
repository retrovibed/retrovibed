import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/library.dart' as lib;
import 'package:retrovibed/media.dart' as media;
import 'api.dart';
import 'socials.row.dart';
import 'socials.details.dart';

// Lists the account's communities, each with Photo/Video/Library/Info
// buttons; Info expands the set of publish plugins the community publishes
// through, which is where they are attached and detached.
class SocialHome extends StatefulWidget {
  final ValueNotifier<media.SearchMode> mode;
  final void Function(media.SearchMode) onModeChanged;
  final FnCommunitySearch search;
  final FnSocialsSearch details;
  final FnSocialsEnable enable;
  final FnSocialsDisable disable;

  const SocialHome({
    super.key,
    required this.mode,
    required this.onModeChanged,
    this.search = communities.search,
    this.details = socials.search,
    this.enable = socials.enable,
    this.disable = socials.disable,
  });

  @override
  State<SocialHome> createState() => _SocialHomeState();
}

class _SocialHomeState extends State<SocialHome> with ds.LoadingState {
  CommunitySearchResponse _resp = CommunitySearchResponse(
    next: CommunitySearchRequest(
      offset: ds.Int64(0),
      limit: ds.Int64(20),
    ),
  );
  Widget _focused = ds.Empty;

  Future<void> _refresh() {
    setState(() => loading = true);
    return httpx
        .withRetry(() => widget.search(_resp.next, options: [authn.request(authn.AuthzCache.meta(context))]))
        .then((response) {
          setState(() {
            _resp = response;
            cause = ds.Error.zero;
          });
        })
        .catchError((cause) {
          setState(() => this.cause = ds.Errors.httpauto(cause, onTap: reseterr));
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((cause) {
          setState(() => this.cause = ds.Error.unknown(cause, onTap: reseterr));
        })
        .whenComplete(() => setState(() => loading = false));
  }

  @override
  void initState() {
    super.initState();
    ds.postframe(_refresh);
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    return Column(
      verticalDirection: defaults.isCompact ? VerticalDirection.up : VerticalDirection.down,
      children: [
        ds.SearchTray(
          autoscroll: true,
          autofocus: defaults.desktop,
          decoration: const InputDecoration(hintText: "search communities"),
          onSubmitted: (v) {
            setState(() {
              _resp.next
                ..query = v
                ..offset = ds.Int64(0);
            });
            return _refresh();
          },
          next: (i) {
            setState(() {
              _resp.next.offset = i;
            });
            _refresh();
          },
          current: _resp.next.offset,
          empty: ds.Int64(_resp.items.length) < _resp.next.limit,
          leading: [
            ds.CompactingMenu.pinned(
              lib.DropdownUpload(
                icon: const Icon(Icons.share),
                help: ds.Hint(const Text("switch to library, discover, or downloads mode")),
                items: [
                  media.SearchModeToggle(
                    mode: media.SearchMode.library,
                    current: widget.mode,
                    icon: Icons.video_library,
                    label: "Library",
                    onSelect: widget.onModeChanged,
                  ),
                  media.SearchModeToggle(
                    mode: media.SearchMode.discovery,
                    current: widget.mode,
                    icon: Icons.travel_explore,
                    label: "Discover",
                    onSelect: widget.onModeChanged,
                  ),
                  media.SearchModeToggle(
                    mode: media.SearchMode.downloads,
                    current: widget.mode,
                    icon: Icons.download,
                    label: "Downloads",
                    onSelect: widget.onModeChanged,
                  ),
                ],
              ),
            ),
          ],
          help: ds.Hint(const Text("search for communities to publish to")),
        ),
        _focused,
        Expanded(
          child: ds.Grid<Community>(
            (context, v) => SocialCommunityRow(
              community: v,
              focused: ValueKey(v.id) == _focused.key,
              onInfo: () => setState(() {
                final key = ValueKey(v.id);
                _focused = key == _focused.key
                    ? ds.Empty
                    : ds.Container(
                        key: key,
                        padding: defaults.padding.copyWith(top: 0, bottom: 0),
                        SocialCommunityDetails(
                          v,
                          search: widget.details,
                          enable: widget.enable,
                          disable: widget.disable,
                        ),
                      );
              }),
            ),
            children: _resp.items,
            loading: loading,
            cause: cause,
            aspectRatio: 3 / 2,
            maxCrossAxisExtent: 420,
            empty: const Center(child: Text('No communities found')),
          ),
        ),
      ],
    );
  }
}
