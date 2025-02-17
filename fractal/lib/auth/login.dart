import 'package:flutter/material.dart';
import 'package:fractal/designkit.dart' as designkit;

class Login extends StatefulWidget {
  final Widget child;
  const Login({required this.child, super.key});

  @override
  State<Login> createState() => _AuthenticatedState();
}

class _AuthenticatedState extends State<Login> {
  bool authenticated = false;
  @override
  Widget build(BuildContext context) {
    return designkit.Full(
      child: designkit.Loading(
        loading: authenticated,
        child: widget.child,
      ),
    );
  }
}
