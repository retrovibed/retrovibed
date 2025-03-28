import 'dart:io';
import 'package:flutter/material.dart';
import 'package:media_kit/media_kit.dart'; // Provides [Player], [Media], [Playlist] etc.
import 'package:console/downloads.dart' as downloads;
import 'package:console/settings.dart' as settings;
import 'package:console/media.dart' as media;
import 'package:console/library.dart' as medialib;
import 'package:console/designkit.dart' as ds;
import 'package:console/design.kit/theme.defaults.dart' as theming;
import 'package:console/mdns.dart' as mdns;
import 'package:console/retrovibed.dart' as retro;

class MyHttpOverrides extends HttpOverrides {
  final List<String> ips;
  MyHttpOverrides({this.ips = const []}) {}
  @override
  HttpClient createHttpClient(SecurityContext? context) {
    return super.createHttpClient(context)
      ..badCertificateCallback = (X509Certificate cert, String host, int port) {
        return ips.any((v) => host == v) ||
            host == "localhost" ||
            host == Platform.localHostname;
      };
  }
}

void main() {
  print("DERP 0 ${retro.bearer_token()}");
  print("DERP 1 ${retro.public_key()}");

  WidgetsFlutterBinding.ensureInitialized();
  MediaKit.ensureInitialized();
  // must come before the daemon since it opens and closes the database.
  HttpOverrides.global = MyHttpOverrides(ips: retro.ips());

  retro.daemon();

  // HttpOverrides.global = MyHttpOverrides();
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
          mdns.MDNSDiscovery(
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
