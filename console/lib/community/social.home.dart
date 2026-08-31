import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/media.dart' as media;
import 'api.dart';

// Lists the account's communities with a toggle per catalog publisher
// (YouTube/Spotify/Instagram/X/etc), so a community owner can pick which
// platforms it publishes to. Actually publishing is not wired up yet.
class SocialHome extends StatefulWidget {
  final ValueNotifier<media.SearchMode> mode;
  final void Function(media.SearchMode) onModeChanged;
  final FnSocialsSearch search;
  final FnSocialsEnable enable;
  final FnSocialsDisable disable;

  const SocialHome({
    super.key,
    required this.mode,
    required this.onModeChanged,
    this.search = socials.search,
    this.enable = socials.enable,
    this.disable = socials.disable,
  });

  @override
  State<SocialHome> createState() => _SocialHomeState();
}

class _SocialHomeState extends State<SocialHome> {
  SocialsSearchResponse _resp = SocialsSearchResponse();
  bool _loading = true;
  Widget _cause = ds.Error.zero;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _resetCause() {
    setState(() => _cause = ds.Error.zero);
  }

  Future<void> _refresh() {
    setState(() => _loading = true);
    return httpx
        .withRetry(() => widget.search(options: [authn.request(authn.AuthzCache.meta(context))]))
        .then((response) {
          setState(() {
            _resp = response;
            _cause = ds.Error.zero;
          });
        })
        .catchError((cause) {
          setState(() => _cause = ds.Errors.httpauto(cause, onTap: _resetCause));
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((cause) {
          setState(() => _cause = ds.Error.unknown(cause, onTap: _resetCause));
        })
        .whenComplete(() => setState(() => _loading = false));
  }

  @override
  void initState() {
    super.initState();
    ds.postframe(_refresh);
  }

  @override
  Widget build(BuildContext context) {
    return ds.Table<CommunitySocial>(
      loading: _loading,
      cause: _cause,
      children: _resp.items,
      empty: const Center(child: Text('No communities found')),
      leading: ds.CompactingMenu.pinned(
        ds.Help(
          PopupMenuButton<String>(
            position: PopupMenuPosition.under,
            color: Theme.of(context).colorScheme.surface,
            surfaceTintColor: Theme.of(context).colorScheme.surface,
            icon: const Icon(Icons.share),
            itemBuilder: (context) => [
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
          ds.Hint(const Text("switch to library, discover, or downloads mode")),
        ),
      ),
      ds.Table.expanded<CommunitySocial>(
        (v) => SocialCommunityRow(
          social: v,
          catalog: _resp.catalog,
          enable: widget.enable,
          disable: widget.disable,
          onChanged: _refresh,
        ),
      ),
    );
  }
}

class SocialCommunityRow extends StatelessWidget {
  final CommunitySocial social;
  final List<PluginPublisher> catalog;
  final FnSocialsEnable enable;
  final FnSocialsDisable disable;
  final VoidCallback onChanged;

  const SocialCommunityRow({
    super.key,
    required this.social,
    required this.catalog,
    required this.enable,
    required this.disable,
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    final enabledIds = social.enabled.map((e) => e.publisherId).toSet();

    return ds.Container(
      padding: defaults.padding,
      Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            social.community.description.isNotEmpty ? social.community.description : social.community.url,
            style: Theme.of(context).textTheme.titleSmall,
          ),
          const SizedBox(height: 4),
          Wrap(
            spacing: 8,
            runSpacing: 4,
            children: catalog.map((p) {
              final enabled = enabledIds.contains(p.id);
              return FilterChip(
                label: Text(p.description.isNotEmpty ? p.description : p.mimetype),
                selected: enabled,
                onSelected: (v) {
                  final auth = [authn.request(authn.AuthzCache.meta(context))];
                  final fut = v
                      ? enable(social.community.id, p.id, options: auth)
                      : disable(social.community.id, p.id, options: auth);
                  httpx.withRetry(() => fut).then((_) => onChanged());
                },
              );
            }).toList(),
          ),
        ],
      ),
    );
  }
}
