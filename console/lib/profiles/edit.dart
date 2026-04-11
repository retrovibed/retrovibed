import 'package:flutter/material.dart';
import 'package:retrovibed/meta.dart' as meta;

class Edit extends StatelessWidget {
  final meta.Profile current;
  final String? pkey;
  final Function(meta.Profile, String key)? onChange;

  Edit(this.current, {super.key, this.onChange, this.pkey = null});

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        TextFormField(
          decoration: InputDecoration(helperText: "name"),
          initialValue: current.display,
          maxLines: 1,
          onChanged: (v) => onChange?.call(current..display = v, pkey ?? ''),
        ),
        if (pkey != null)
          TextFormField(
            decoration: InputDecoration(helperText: "public key"),
            initialValue: pkey,
            maxLines: 1,
            onChanged: (v) => onChange?.call(current, v),
          ),
      ],
    );
  }
}
