import 'dart:io';

import 'package:flutter/material.dart';
import 'package:window_manager/window_manager.dart';

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
import 'package:retrovibed/design.kit/theme.defaults.dart' as theming;
import 'package:retrovibed/design.kit/modals.dart' as modals;
import 'package:retrovibed/community.dart' as community;
import 'package:retrovibed/mimex.dart' as mimex;

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

void main(List<String> args) async {
  // --smoke: initialize exactly as a normal launch would (proving things like
  // media_kit/libmpv wiring loaded successfully), then exit gracefully
  // instead of presenting the UI. Used by the AppImage smoke test to verify
  // the shipped artifact starts cleanly without needing a real display session.
  final smoke = args.contains('--smoke');

  await retro.run(() {
    if (smoke) {
      exit(0);
    }

    FlutterError.onError = FlutterError.dumpErrorToConsole;

    ErrorWidget.builder = (FlutterErrorDetails details) {
      FlutterError.dumpErrorToConsole(details);
      return Material(
        child: InkWell(
          onTap: () => ds.postframe(() => runApp(Retrovibed())),
          child: const Center(child: Text('Something went wrong. Tap to restart.')),
        ),
      );
    };

    runApp(Retrovibed());
  });
}

final ds.AsyncVoidCallback _startdaemon = ds.toasync(ds.once(retro.daemon));

class Retrovibed extends StatelessWidget {
  Retrovibed({super.key});

  @override
  Widget build(BuildContext context) {
    final defaults = theming.Defaults.defaults;
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
        return MediaQuery(
          data: MediaQuery.of(context).copyWith(textScaler: autoscaling(context)),
          child: Theme(
            data: Theme.of(context).copyWith(extensions: [defaults.copyWith(isCompact: isCompact)]),
            child: child ?? ds.Empty,
          ),
        );
      },
      home: Material(
        child: SafeArea(
          child: ds.Full(
            ds.LoadingGuard(
              ds.ErrorBoundary(
                ds.HelpScope(
                  authn.Login(
                    authenticated: _startdaemon,
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
                                            mimetype: mimex.category(q.mimetypes),
                                          ),
                                          options: [authn.request(authn.AuthzCache.meta(ctx))],
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
                                      if (defaults.desktop) {
                                        tabs = DragToMoveArea(child: tabs);
                                      }

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
                                                downloads.AutoHelp(
                                                  downloads.MeteredWarning(const downloads.Display()),
                                                ),
                                              ),
                                              modals.Node(community.AutoHelp(community.Management())),
                                              modals.Node(settings.AutoHelp(const settings.Display())),
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
        ),
      ),
    );
  }
}
