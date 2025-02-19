import 'package:flutter/material.dart';
import 'package:fractal/downloads.dart' as downloads;
import 'package:fractal/discovery.dart' as discovery;
import 'dart:io';

class MyHttpOverrides extends HttpOverrides {
  @override
  HttpClient createHttpClient(SecurityContext? context) {
    return super.createHttpClient(context)
      ..badCertificateCallback = (X509Certificate cert, String host, int port) {
        return host == "localhost";
      };
  }
}

void main() {
  HttpOverrides.global = MyHttpOverrides();
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      darkTheme: ThemeData(brightness: Brightness.dark),
      themeMode: ThemeMode.dark,
      home: DefaultTabController(
        length: 4,
        child: Scaffold(
          appBar: TabBar(
            tabs: [
              Tab(icon: Icon(Icons.share)),
              Tab(icon: Icon(Icons.movie)),
              Tab(icon: Icon(Icons.download)),
              Tab(icon: Icon(Icons.settings)),
            ],
          ),
          body: TabBarView(
            children: [
              discovery.Display(),
              Icon(Icons.movie),
              downloads.List(),
              Icon(Icons.settings),
            ],
          ),
        ),
      ),
    );
  }
}
