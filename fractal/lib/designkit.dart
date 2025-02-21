import 'package:flutter/material.dart';
import 'design.kit/theme.defaults.dart' as theming;
export 'design.kit/screens.dart';
export 'design.kit/accordian.dart';
export 'design.kit/errors.dart';

theming.Defaults theme(BuildContext context) {
  return Theme.of(context).extension<theming.Defaults>() ??
      theming.Defaults.defaults;
}
