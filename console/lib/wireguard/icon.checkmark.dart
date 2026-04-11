import 'package:flutter/material.dart';

class IconCheckmark extends StatelessWidget {
  static Future<void> _noop() {
    return Future.value(null);
  }

  final bool checked;
  final Future<void> Function()? onTap;
  const IconCheckmark(this.checked, {super.key, this.onTap = _noop});

  @override
  Widget build(BuildContext context) {
    return IconButton(
      hoverColor: onTap == _noop ? Colors.transparent : null,
      onPressed: onTap,
      icon: Opacity(
        opacity: checked ? 1.0 : 0.1,
        child: Icon(Icons.check, color: Colors.lightGreenAccent),
      ),
    );
  }
}
