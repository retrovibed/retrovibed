import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';

class Hyperlink extends StatelessWidget {
  final String text;
  final Uri uri;
  final TextStyle? style;

  const Hyperlink(
    this.text, {
    super.key,
    required this.uri,
    this.style,
  });

  Hyperlink.fromString(
    this.text, {
    super.key,
    required String url,
    this.style,
  }) : uri = Uri.parse(url);

  static WidgetSpan inline(
    String text, {
    required String url,
    TextStyle? style,
    Key? key,
  }) {
    return WidgetSpan(
      alignment: PlaceholderAlignment.baseline,
      baseline: TextBaseline.alphabetic,
      child: Hyperlink.fromString(text, url: url, style: style, key: key),
    );
  }

  @override
  Widget build(BuildContext context) {
    final base = style ?? Theme.of(context).textTheme.bodyMedium;
    return GestureDetector(
      onTap: () => launchUrl(uri),
      child: Text(
        text,
        style: base?.copyWith(
          decoration: TextDecoration.underline,
          color: Theme.of(context).colorScheme.primary,
        ),
      ),
    );
  }
}
