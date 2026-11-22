/// Dates as the site shows them: "March 2027" from a frontmatter date such as "2027-03"
/// or "2027-03-14". Anything else is shown as written.
String formatMonth(String date) {
  const months = [
    'January', 'February', 'March', 'April', 'May', 'June',
    'July', 'August', 'September', 'October', 'November', 'December',
  ];
  final parts = date.split('-');
  final month = parts.length >= 2 ? int.tryParse(parts[1]) : null;
  if (month == null || month < 1 || month > 12) return date;
  return '${months[month - 1]} ${parts[0]}';
}
