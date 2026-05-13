import 'package:flutter/material.dart';
import 'help.scope.dart';

class HelpGlobal extends StatefulWidget {
  final Widget child;
  final Widget description;
  const HelpGlobal(this.child, this.description, {super.key});

  @override
  State<HelpGlobal> createState() => _HelpGlobalState();
}

class _HelpGlobalState extends State<HelpGlobal> {
  HelpScopeState? _scope;

  @override
  void initState() {
    super.initState();
    _scope = HelpScope.of(context);
    _scope?.registerGlobal(widget.description);
  }

  @override
  void dispose() {
    _scope?.unregisterGlobal(widget.description);
    super.dispose();
  }

  @override
  void didUpdateWidget(HelpGlobal old) {
    super.didUpdateWidget(old);
    if (old.description != widget.description) {
      _scope?.unregisterGlobal(old.description);
      _scope?.registerGlobal(widget.description);
    }
  }

  @override
  Widget build(BuildContext context) => widget.child;
}
