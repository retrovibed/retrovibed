import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/profiles.dart' as profiles;

class Create extends StatelessWidget {
  final meta.Profile profile;
  final String publicKey;
  final meta.Token token;
  final Function(meta.Profile, String, meta.Token)? onChange;

  const Create(
    this.profile,
    this.publicKey,
    this.token, {
    super.key,
    this.onChange,
  });

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return LayoutBuilder(
      builder: (context, constraints) {
        const double minItemWidth = 300;
        final spacing = defaults.spacing;

        // Calculate width for 2 items per row
        double itemWidth = (constraints.maxWidth - spacing) / 2;

        // If equal width is smaller than the minimum, let them take full width (stacking)
        if (itemWidth < minItemWidth) {
          itemWidth = constraints.maxWidth;
        }

        return Wrap(
          spacing: spacing,
          runSpacing: spacing,
          children: [
            SizedBox(
              width: itemWidth,
              child: profiles.Edit(
                profile,
                pkey: publicKey,
                onChange: (p, key) => onChange?.call(p, key, token),
              ),
            ),
            SizedBox(
              width: itemWidth,
              child: profiles.AuthzMetaEdit(
                token,
                onChange: (t) => onChange?.call(profile, publicKey, t),
              ),
            ),
          ],
        );
      },
    );
  }
}
