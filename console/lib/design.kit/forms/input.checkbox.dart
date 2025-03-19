import 'package:flutter/material.dart' as m;
import 'package:flutter/material.dart';

class Checkbox extends StatelessWidget {
  final bool value;
  final void Function(bool?)? onChanged;
  Checkbox({super.key, this.value = false, this.onChanged});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [m.Checkbox(value: value, onChanged: onChanged), Spacer()],
    );
  }
}
