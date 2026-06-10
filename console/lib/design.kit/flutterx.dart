import 'package:flutter/material.dart';

void postframe(VoidCallback fn) {
  WidgetsBinding.instance.addPostFrameCallback((_) => fn());
}
