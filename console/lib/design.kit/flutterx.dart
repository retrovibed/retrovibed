import 'package:material_ui/material_ui.dart';

void postframe(VoidCallback fn) {
  WidgetsBinding.instance.addPostFrameCallback((_) => fn());
}
