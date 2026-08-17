import 'package:flutter/material.dart';
import 'package:window_manager/window_manager.dart';

import 'package:retrovibed/navbar.dart' as navbar;
import 'package:retrovibed/remote.dart' as remote;
import 'package:retrovibed/settings.dart' as settings;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/library.dart' as medialib;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/modals.dart' as modals;
import 'package:retrovibed/community.dart' as community;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/meta.dart' as meta;

/// The app's tab-based navigation: movie/library, remote control, community,
/// and settings, each behind a [DefaultTabController]-driven [TabBar] and
/// [TabBarView]. Self-contained so it can be tested without the auth/media
/// playback/network gates that wrap it in main.dart.
class Routes extends StatelessWidget {
  final Future<meta.DaemonSearchResponse> Function(meta.DaemonSearchRequest) daemonSearch;
  final Future<Stream<meta.Daemon>> Function({List<httpx.Option> options}) daemonDiscover;

  const Routes({
    super.key,
    this.daemonSearch = meta.daemons.search,
    this.daemonDiscover = meta.daemons.discover,
  });

  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
      length: 4,
      child: ds.build((context) {
        final defaults = ds.Defaults.of(context);
        final compact = defaults.isCompact;
        final nochrome = ds.Full.nochrome(context);
        final tabbar = TabBar(
          dividerHeight: 0,
          tabs: [
            Tab(icon: Icon(Icons.movie)),
            Tab(icon: Icon(Icons.settings_remote)),
            Tab(icon: Icon(Icons.groups)),
            Tab(icon: Icon(Icons.settings)),
          ],
        );

        Widget tabs = DecoratedBox(
          decoration: BoxDecoration(
            border: Border(
              bottom: BorderSide(
                width: 1.0,
                color: Theme.of(context).dividerColor,
              ),
            ),
          ),
          child: DragToMoveArea(
            child: Row(
              children: [
                Expanded(
                  child: tabbar,
                ),
                if (defaults.desktop) navbar.Hamburger(),
              ],
            ),
          ),
        );

        return Scaffold(
          appBar: (!compact && !nochrome)
              ? PreferredSize(
                  preferredSize: Size.fromHeight(kTextTabBarHeight),
                  child: tabs,
                )
              : null,
          bottomNavigationBar: (compact && !nochrome) ? tabs : null,
          body: ds.ErrorBoundary(
            TabBarView(
              // disable scrolling so that people dont accidently scroll
              // through the tabs
              physics: const NeverScrollableScrollPhysics(),
              children: [
                modals.Node(
                  media.AutoHelp(
                    media.Playlist.wrap((ctx, s) {
                      return media.VideoScreen(
                        medialib.Home(
                          key: ValueKey("library"),
                          focus: defaults.mobile ? null : s.searchfocus,
                          controller: s.controller,
                          highlighted: s.known.id,
                          search: s.search,
                        ),
                        s.player,
                        s.playerfocus,
                        s.overlay,
                      );
                    }),
                  ),
                ),
                modals.Node(
                  media.Playlist.wrap((ctx, s) {
                    return remote.Connect(search: s.search, daemonDiscover: daemonDiscover);
                  }),
                ),
                modals.Node(community.AutoHelp(community.Management())),
                modals.Node(
                  settings.AutoHelp(
                    settings.Display(daemonSearch: daemonSearch, daemonDiscover: daemonDiscover),
                  ),
                ),
              ],
            ),
          ),
        );
      }),
    );
  }
}
