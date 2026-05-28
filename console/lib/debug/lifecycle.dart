import 'package:flutter/material.dart';

class Lifecycle extends StatefulWidget {
  final Widget child;
  final String message;

  const Lifecycle(this.child, {super.key, this.message = 'lifecycle detection'});

  @override
  State<StatefulWidget> createState() => _Lifecycle();
}

class _Lifecycle extends State<Lifecycle> with WidgetsBindingObserver {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    super.dispose();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    print("${widget.message}: $state");
  }

  @override
  Widget build(BuildContext context) => widget.child;
}
