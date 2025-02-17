import 'package:flutter/material.dart';
import 'package:fractal/downloads.dart' as downloads;
import 'package:fractal/discovery.dart' as discovery;

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return const MaterialApp(
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
