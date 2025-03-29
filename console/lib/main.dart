import 'package:flutter/material.dart';
import 'package:media_kit/media_kit.dart'; // Provides [Player], [Media], [Playlist] etc.
import 'package:console/downloads.dart' as downloads;
import 'package:console/settings.dart' as settings;
import 'package:console/media.dart' as media;
import 'package:console/library.dart' as medialib;
import 'package:console/designkit.dart' as ds;
import 'package:console/meta.dart' as meta;
import 'package:console/retrovibed.dart' as retro;
import 'package:console/design.kit/theme.defaults.dart' as theming;

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  MediaKit.ensureInitialized();

  retro.daemon();

  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      darkTheme: ThemeData(
        brightness: Brightness.dark,
        cardTheme: CardTheme(margin: EdgeInsets.all(10.0)),
        extensions: [theming.Defaults.defaults],
      ),
      themeMode: ThemeMode.dark,
      home: Material(
        child: ds.Full(
          meta.DaemonAuto(
            media.Playlist(
              DefaultTabController(
                length: 3,
                child: Scaffold(
                  appBar: TabBar(
                    tabs: [
                      Tab(icon: Icon(Icons.movie)),
                      Tab(icon: Icon(Icons.download)),
                      Tab(icon: Icon(Icons.settings)),
                    ],
                  ),
                  body: TabBarView(
                    children: [
                      ds.ErrorBoundary(
                        media.Playlist.wrap((ctx, s) {
                          return media.VideoScreen(
                            medialib.AvailableListDisplay(
                              focus: s.searchfocus,
                              controller: s.controller,
                            ),
                            s.player,
                          );
                        }),
                      ),
                      ds.ErrorBoundary(downloads.Display()),
                      ds.ErrorBoundary(settings.Display()),
                    ],
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
