import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/library.dart' as lib;
import 'package:retrovibed/media.dart' as media;
import 'api.dart';
import 'socials.card.dart';
import 'socials.details.dart';

// Lists the account's communities, each with Photo/Video/Library/Info
// buttons; Info expands the set of publish plugins the community publishes
// through, which is where they are attached and detached.
class SocialHome extends StatefulWidget {
  final ValueNotifier<media.SearchMode> mode;
  final void Function(media.SearchMode) onModeChanged;
  final ValueNotifier<media.MediaSearchState> search;
  final FnCommunitySearch apisearch;
  final FnSocialsSearch apidetails;
  final FnSocialsEnable apienable;
  final FnSocialsDisable apidisable;

  const SocialHome({
    super.key,
    required this.mode,
    required this.onModeChanged,
    required this.search,
    this.apisearch = communities.search,
    this.apidetails = socials.search,
    this.apienable = socials.enable,
    this.apidisable = socials.disable,
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
        .withRetry(() => widget.apisearch(_resp.next, options: [authn.request(authn.AuthzCache.meta(context))]))
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
              lib.DropdownNavMenu.options(
                icon: const Icon(Icons.share),
                search: widget.search,
                mode: widget.mode,
                onModeChanged: widget.onModeChanged,
              ),
            ),
          ],
          help: ds.Hint(const Text("search for communities to publish to")),
        ),
        Expanded(
          child: ds.Grid<Community>(
            (context, v) => SocialCard(
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
                          search: widget.apidetails,
                          enable: widget.apienable,
                          disable: widget.apidisable,
                        ),
                      );
              }),
            ),
            leading: [_focused],
            physics: const AlwaysScrollableScrollPhysics(),
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
