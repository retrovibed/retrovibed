import 'dart:async';

import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:flutter/material.dart';
import 'package:retrovibed/retrovibed.dart' as retro;

bool _isMetered(List<ConnectivityResult> results) {
  return results.any((r) => r == ConnectivityResult.mobile);
}

class Metered extends StatefulWidget {
  final Widget child;

  const Metered(this.child, {super.key});

  @override
  State<Metered> createState() => _MeteredState();
}

class _MeteredState extends State<Metered> {
  StreamSubscription<List<ConnectivityResult>>? _sub;

  @override
  void initState() {
    super.initState();
    _sub = Connectivity().onConnectivityChanged.listen((results) {
      print("metered ${results}");
      retro.set_metered(_isMetered(results));
    });
    Connectivity().checkConnectivity().then((results) {
      retro.set_metered(_isMetered(results));
    }).ignore();
  }

  @override
  void dispose() {
    _sub?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) => widget.child;
}
