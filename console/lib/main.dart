import 'dart:io';

import 'package:flutter/material.dart';
import 'package:media_kit/media_kit.dart'; // Provides [Player], [Media], [Playlist] etc.

import 'package:retrovibed/navbar.dart' as navbar;
import 'package:retrovibed/billing.dart' as billing;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/downloads.dart' as downloads;
import 'package:retrovibed/settings.dart' as settings;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/library.dart' as medialib;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/retrovibed.dart' as retro;
import 'package:retrovibed/deeplink.dart';
import 'package:retrovibed/env.dart' as env;
import 'package:retrovibed/design.kit/theme.defaults.dart' as theming;
import 'package:retrovibed/design.kit/modals.dart' as modals;
import 'package:retrovibed/community.dart' as community;
import 'package:window_manager/window_manager.dart';

TextScaler autoscaling(BuildContext context) {
  // final width = MediaQuery.of(context).size.width;
  // print("autoscaling width ${width} - ${MediaQuery.of(context).size.height}");
  // if (width > 1920) {
  //   return TextScaler.linear(
  //     2.0,
  //   ).clamp(minScaleFactor: 0.8, maxScaleFactor: 4.0);
  // } else {
  //   return TextScaler.linear(1.0);
  // }

  return TextScaler.linear(1.0);
}

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  HttpOverrides.global = meta.DaemonHttpOverrides();
  retro.logging();
  await env.xdg();

  if (theming.Defaults.defaults.desktop) {
    await windowManager.ensureInitialized();
    await Future.wait([
      windowManager.setTitleBarStyle(TitleBarStyle.hidden),
      windowManager.maximize(),
    ]);
  }

  MediaKit.ensureInitialized();
  runApp(Retrovibed());
}

class Retrovibed extends StatelessWidget {
  Retrovibed({super.key});

  @override
  Widget build(BuildContext context) {
    final btnstyle = ButtonStyle(
      mouseCursor: WidgetStateProperty.all(SystemMouseCursors.click),
    );
    return MaterialApp(
      debugShowCheckedModeBanner: false,
      darkTheme: ThemeData(
        brightness: Brightness.dark,
        hoverColor: ds.Defaults.defaults.highlight,
        iconButtonTheme: IconButtonThemeData(
          style: btnstyle,
        ),
        textButtonTheme: TextButtonThemeData(
          style: btnstyle,
        ),
        outlinedButtonTheme: OutlinedButtonThemeData(
          style: btnstyle,
        ),
        elevatedButtonTheme: ElevatedButtonThemeData(
          style: btnstyle,
        ),
        filledButtonTheme: FilledButtonThemeData(
          style: btnstyle,
        ),
        segmentedButtonTheme: SegmentedButtonThemeData(
          style: btnstyle,
        ),
        cardTheme: CardThemeData(margin: EdgeInsets.all(10.0)),
        inputDecorationTheme: InputDecorationTheme(
          focusColor: Colors.transparent,
        ),
        popupMenuTheme: PopupMenuThemeData(),
      ),
      themeMode: ThemeMode.dark,
      builder: (context, child) {
        final isCompact =
            MediaQuery.of(context).size.width < theming.Defaults.defaults.compact || theming.Defaults.defaults.mobile;
        final defaults = theming.Defaults.defaults.copyWith(isCompact: isCompact);
        return MediaQuery(
          data: MediaQuery.of(context).copyWith(textScaler: autoscaling(context)),
          child: Theme(
            data: Theme.of(context).copyWith(extensions: [defaults]),
            child: child ?? ds.Empty,
          ),
        );
      },
      home: Material(
        child: SafeArea(
          child: ds.Full(
            ds.HelpScope(
              authn.Login(
                authenticated: () async {
                  retro.daemon();
                },
                meta.EndpointAuto(
                  authn.Authenticated(
                    authn.DeeppoolAuthzCache(
                      authn.AuthzCache(
                        DeepLink(
                          billing.Registered(
                            media.Playlist(
                              tracing: (ctx, pos, dur, q, id) {
                                medialib.recent
                                    .record(
                                      medialib.RecentRecordRequest(
                                        media: medialib.Media(id: id),
                                        position: ds.Int64(pos.inMilliseconds),
                                        duration: ds.Int64(dur.inMilliseconds),
                                        query: q,
                                      ),
                                      options: [authn.AuthzCache.bearer(ctx)],
                                    )
                                    .then((v) {})
                                    .catchError((cause) {
                                      print(
                                        "failed to record watch event ${pos}/${dur} - ${q} - ${cause}",
                                      );
                                    })
                                    .ignore();
                              },
                              DefaultTabController(
                                length: 4,
                                child: ds.build((context) {
                                  final defaults = ds.Defaults.of(context);
                                  final compact = defaults.isCompact;
                                  final nochrome = ds.Full.nochrome(context);
                                  final tabbar = TabBar(
                                    dividerHeight: 0,
                                    tabs: [
                                      Tab(icon: Icon(Icons.movie)),
                                      Tab(icon: Icon(Icons.download)),
                                      Tab(icon: Icon(Icons.groups)),
                                      Tab(icon: Icon(Icons.settings)),
                                    ],
                                  );

                                  final tabs = DecoratedBox(
                                    decoration: BoxDecoration(
                                      border: Border(
                                        bottom: BorderSide(
                                          width: 1.0,
                                          color: Theme.of(context).dividerColor,
                                        ),
                                      ),
                                    ),
                                    child: Row(
                                      children: [
                                        Expanded(child: tabbar),
                                        if (defaults.desktop) navbar.Hamburger(),
                                      ],
                                    ),
                                  );
                                  return Scaffold(
                                    appBar:
                                        (!compact && !nochrome)
                                            ? PreferredSize(
                                              preferredSize: Size.fromHeight(kTextTabBarHeight),
                                              child: tabs,
                                            )
                                            : null,
                                    bottomNavigationBar: (compact && !nochrome) ? tabs : null,
                                    body: ds.ErrorBoundary(
                                      TabBarView(
                                        children: [
                                          modals.Node(
                                            media.Playlist.wrap((ctx, s) {
                                              return media.VideoScreen(
                                                medialib.AvailableGridDisplay(
                                                  focus: defaults.mobile ? null : s.searchfocus,
                                                  controller: s.controller,
                                                  highlighted: s.current.id,
                                                  search: s.search,
                                                ),
                                                s.player,
                                                s.playerfocus,
                                              );
                                            }),
                                          ),
                                          modals.Node(downloads.Display()),
                                          modals.Node(community.Management()),
                                          modals.Node(settings.Display()),
                                        ],
                                      ),
                                    ),
                                  );
                                }),
                              ),
                            ),
                          ),
                        ),
                      ),
                    ),
                  ),
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}
