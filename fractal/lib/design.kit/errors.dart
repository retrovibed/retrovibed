import 'package:flutter/material.dart';

class Error extends StatelessWidget {
  final Widget child;

  const Error({super.key, required this.child});

  static Error text(String text) => Error(child: Text(text));
  static Error unknown(Object obj) {
    print("${obj.toString()}");
    return Error(child: Text("an unexpected problem has occurred"));
  }

  @override
  Widget build(BuildContext context) {
    return Container(child: Center(child: child));
  }
}
