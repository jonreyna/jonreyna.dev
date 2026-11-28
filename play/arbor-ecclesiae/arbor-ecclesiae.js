/* Arbor Ecclesiae: the branching of the church, c. 49 to 1700, as a radial tree.
   Needs D3 v7 (d3.min.js, loaded first). The page's markup is arbor-ecclesiae.html
   and its styles arbor-ecclesiae.css; see play/README.md. */
(function () {
"use strict";
if (!window.d3) { document.getElementById("nod3").hidden = false; return; }

/* =====================================================================
   1. DATA
   Each node: id, name, date (year AD), doctrine, size (influence 1–10),
   extinct, uncertain, reference (Wikipedia URL or null), children.
   Optional: end (extinction year), circa / endCirca (show "c."),
   tradition (colour family, inherited when absent), event (council or
   schism knot on the trunk), eventKind, label (short display name).
   Children are listed in clockwise order.

   Influence rubric: 9–10 dominant, continent-scale communions and the
   councils they rest on; 6–8 major lasting families or very large
   movements; 3–5 significant regional or long-lived movements;
   1–2 small, short-lived or individual.
   ===================================================================== */
const wiki = slug => "https://en.wikipedia.org/wiki/" + slug;

const DATA = {
  id: "apostolic-church", name: "Apostolic Church — Council of Jerusalem", label: "Apostolic Church",
  date: 49, circa: true, size: 10, extinct: false, uncertain: false,
  tradition: "trunk", event: true, eventKind: "Council",
  reference: wiki("Council_of_Jerusalem"),
  doctrine: "The first-generation church of the apostles. At Jerusalem the apostles and elders ruled that Gentile converts need not be circumcised or keep the full Law of Moses (Acts 15).",
  children: [
    { id: "judaizers", name: "Judaizers", date: 49, circa: true, end: 135, endCirca: true, size: 3, extinct: true, uncertain: false,
      reference: wiki("Judaizers"),
      doctrine: "Jewish Christians who held that Gentile believers must be circumcised and keep the Law of Moses, the position rejected by the Jerusalem Council and Paul’s letters.",
      children: [
        { id: "ebionites", name: "Ebionites", date: 70, circa: true, end: 400, endCirca: true, size: 3, extinct: true, uncertain: false,
          reference: wiki("Ebionites"),
          doctrine: "A Law-observant Jewish-Christian sect that accepted Jesus as Messiah but saw him as a human prophet, not divine, and rejected Paul’s letters." },
        { id: "nazarenes", name: "Nazarenes", date: 70, circa: true, end: 400, endCirca: true, size: 2, extinct: true, uncertain: true,
          reference: wiki("Nazarene_(sect)"),
          doctrine: "Torah-observant Jewish Christians described by Epiphanius and Jerome. Unlike the Ebionites they seem to have accepted Christ’s divinity, but their identity and continuity are debated." }
      ] },
    { id: "gnosticism", name: "Gnosticism", date: 80, circa: true, end: 500, endCirca: true, size: 6, extinct: true, uncertain: false,
      reference: wiki("Gnosticism"),
      doctrine: "A family of movements teaching salvation through secret knowledge (gnosis), a flawed creator distinct from the supreme God, and a material world to be escaped.",
      children: [
        { id: "docetism", name: "Docetism", date: 100, circa: true, end: 250, endCirca: true, size: 3, extinct: true, uncertain: false,
          reference: wiki("Docetism"),
          doctrine: "The view that Christ only seemed (dokein) to have a body and to suffer. Ignatius of Antioch wrote against it." },
        { id: "basilideans", name: "Basilideans", date: 130, circa: true, end: 350, endCirca: true, size: 3, extinct: true, uncertain: false,
          reference: wiki("Basilides"),
          doctrine: "Followers of Basilides of Alexandria, who taught an elaborate cosmology of heavens and powers descending from an unknowable God." },
        { id: "valentinianism", name: "Valentinianism", date: 140, circa: true, end: 450, endCirca: true, size: 5, extinct: true, uncertain: false,
          reference: wiki("Valentinianism"),
          doctrine: "The school of Valentinus, the most influential Gnostic movement. Its myth of divine aeons and the fall of Sophia was read alongside the church’s own scriptures." },
        { id: "sethianism", name: "Sethianism", date: 150, circa: true, end: 400, endCirca: true, size: 3, extinct: true, uncertain: true,
          reference: wiki("Sethianism"),
          doctrine: "A Gnostic tradition centred on Seth as a bearer of divine knowledge, known mainly from the Nag Hammadi texts. Whether it was ever an organised group is debated." }
      ] },
    { id: "marcionism", name: "Marcionism", date: 144, end: 450, endCirca: true, size: 5, extinct: true, uncertain: false,
      reference: wiki("Marcionism"),
      doctrine: "Marcion of Sinope taught that the God of the Hebrew Bible was not the loving Father revealed by Jesus, and issued a canon of Luke and ten letters of Paul. He was excommunicated at Rome in 144." },
    { id: "montanism", name: "Montanism", date: 156, circa: true, end: 550, endCirca: true, size: 4, extinct: true, uncertain: false,
      reference: wiki("Montanism"),
      doctrine: "The ‘New Prophecy’ of Montanus, Priscilla and Maximilla in Phrygia: ecstatic revelation, strict fasting and moral rigour. It won over Tertullian." },
    { id: "great-church", name: "The Great Church", label: "Great Church", date: 180, circa: true, size: 8, extinct: false, uncertain: false,
      event: true, eventKind: "Era",
      reference: wiki("Great_Church"),
      doctrine: "The broad ‘proto-orthodox’ communion of bishops in apostolic succession, held together by the rule of faith and an emerging New Testament canon, as defended by Irenaeus of Lyon.",
      children: [
        { id: "adoptionism", name: "Adoptionism", date: 190, circa: true, end: 300, endCirca: true, size: 3, extinct: true, uncertain: false,
          reference: wiki("Adoptionism"),
          doctrine: "The teaching that Jesus was a man adopted as Son of God at his baptism or resurrection. Taught by Theodotus of Byzantium; Paul of Samosata was condemned at Antioch in 268." },
        { id: "sabellianism", name: "Sabellianism", date: 215, circa: true, end: 400, endCirca: true, size: 4, extinct: true, uncertain: false,
          reference: wiki("Sabellianism"),
          doctrine: "Modalism: Father, Son and Spirit are three modes or roles of one divine person, not distinct persons. Sabellius was excommunicated at Rome c. 220." },
        { id: "novatianism", name: "Novatianism", date: 251, end: 650, endCirca: true, size: 4, extinct: true, uncertain: false,
          reference: wiki("Novatianism"),
          doctrine: "A rigorist schism led by Novatian of Rome, who refused to readmit Christians who had lapsed in the Decian persecution. Orthodox in doctrine, it lasted for centuries." },
        { id: "donatism", name: "Donatism", date: 312, end: 700, endCirca: true, size: 5, extinct: true, uncertain: false,
          reference: wiki("Donatism"),
          doctrine: "A North African schism that rejected clergy who had handed over scriptures under persecution, holding that a sacrament depends on the minister’s holiness. Opposed by Augustine." },
        { id: "nicaea", name: "First Council of Nicaea", label: "Nicaea 325", date: 325, size: 9, extinct: false, uncertain: false,
          event: true, eventKind: "Ecumenical council",
          reference: wiki("First_Council_of_Nicaea"),
          doctrine: "The first ecumenical council, called by Constantine. It declared the Son ‘of one substance’ (homoousios) with the Father and issued the original Nicene Creed.",
          children: [
            { id: "arianism", name: "Arianism", date: 325, end: 700, endCirca: true, size: 7, extinct: true, uncertain: false,
              reference: wiki("Arianism"),
              doctrine: "Arius taught that the Son was created by the Father: ‘there was when he was not’. Condemned at Nicaea, it still prevailed for decades and among the Germanic peoples.",
              children: [
                { id: "semi-arians", name: "Semi-Arians", date: 358, circa: true, end: 400, endCirca: true, size: 4, extinct: true, uncertain: false,
                  reference: wiki("Semi-Arianism"),
                  doctrine: "The Homoiousians held that the Son is ‘of like substance’ with the Father. Many, led by Basil of Ancyra, were eventually reconciled to the Nicene faith." },
                { id: "anomoeans", name: "Anomoeans", date: 360, circa: true, end: 400, endCirca: true, size: 2, extinct: true, uncertain: false,
                  reference: wiki("Anomoeanism"),
                  doctrine: "Radical Arians led by Aetius and Eunomius, who taught that the Son is ‘unlike’ (anomoios) the Father in essence." },
                { id: "gothic-arianism", name: "Gothic Arianism", date: 341, circa: true, end: 671, endCirca: true, size: 5, extinct: true, uncertain: false,
                  reference: wiki("Gothic_Christianity"),
                  doctrine: "Homoian Christianity brought to the Goths by Ulfilas, who translated the Bible into Gothic. The Visigoths became Nicene in 589 and the Lombards by the late 600s." }
              ] },
            { id: "constantinople-i", name: "First Council of Constantinople", label: "Constantinople 381", date: 381, size: 8, extinct: false, uncertain: false,
              event: true, eventKind: "Ecumenical council",
              reference: wiki("First_Council_of_Constantinople"),
              doctrine: "The second ecumenical council. It affirmed the full divinity of the Holy Spirit and expanded the creed into the Niceno-Constantinopolitan form still used today.",
              children: [
                { id: "apollinarism", name: "Apollinarism", date: 381, end: 420, endCirca: true, size: 2, extinct: true, uncertain: false,
                  reference: wiki("Apollinarism"),
                  doctrine: "Apollinaris of Laodicea taught that in Christ the divine Logos took the place of a human mind. Condemned in 381." },
                { id: "pneumatomachi", name: "Pneumatomachi", date: 381, end: 430, endCirca: true, size: 2, extinct: true, uncertain: false,
                  reference: wiki("Pneumatomachi"),
                  doctrine: "The ‘fighters against the Spirit’, or Macedonians, accepted the Son’s divinity but denied that of the Holy Spirit." },
                { id: "helvidius", name: "Helvidius", date: 383, circa: true, end: 400, endCirca: true, size: 1, extinct: true, uncertain: true,
                  reference: wiki("Helvidius"),
                  doctrine: "Helvidius argued that Mary had children with Joseph after Jesus’ birth, and Jerome answered him c. 383. Later writers call holders of this view Antidicomarianites, but there is little sign of an organised following." },
                { id: "priscillianism", name: "Priscillianism", date: 385, end: 600, endCirca: true, size: 2, extinct: true, uncertain: false,
                  reference: wiki("Priscillianism"),
                  doctrine: "An ascetic movement in Hispania founded by Priscillian, who was accused of Gnostic and Manichaean ideas. His execution in 385 was the first for heresy." },
                { id: "pelagianism", name: "Pelagianism", date: 405, circa: true, end: 450, endCirca: true, size: 5, extinct: true, uncertain: false,
                  reference: wiki("Pelagianism"),
                  doctrine: "Pelagius taught that people can obey God without special grace and denied inherited original sin. Opposed by Augustine, it was condemned at Carthage in 418 and at Ephesus in 431.",
                  children: [
                    { id: "semi-pelagianism", name: "Semi-Pelagianism", date: 425, circa: true, end: 529, size: 3, extinct: true, uncertain: false,
                      reference: wiki("Semi-Pelagianism"),
                      doctrine: "Monks of southern Gaul such as John Cassian held that the first step toward faith comes from the human will, which grace then assists. Rejected at the Council of Orange in 529." }
                  ] },
                { id: "ephesus", name: "Council of Ephesus", label: "Ephesus 431", date: 431, size: 8, extinct: false, uncertain: false,
                  event: true, eventKind: "Ecumenical council",
                  reference: wiki("Council_of_Ephesus"),
                  doctrine: "The third ecumenical council. It deposed Nestorius and affirmed Mary as Theotokos (‘God-bearer’), insisting that Christ is one person.",
                  children: [
                    { id: "chalcedon", name: "Council of Chalcedon", label: "Chalcedon 451", date: 451, size: 9, extinct: false, uncertain: false,
                      event: true, eventKind: "Ecumenical council",
                      reference: wiki("Council_of_Chalcedon"),
                      doctrine: "The fourth ecumenical council. Following Pope Leo’s Tome, it defined Christ as one person ‘in two natures, without confusion, change, division or separation’.",
                      children: [
                        { id: "constantinople-ii", name: "Second Council of Constantinople", label: "Constantinople II 553", date: 553, size: 5, extinct: false, uncertain: false,
                          event: true, eventKind: "Ecumenical council",
                          reference: wiki("Second_Council_of_Constantinople"),
                          doctrine: "The fifth ecumenical council, under Justinian. It condemned the ‘Three Chapters’ in an effort to win back the miaphysites, reading Chalcedon in the light of Cyril.",
                          children: [
                            { id: "nicaea-ii", name: "Second Council of Nicaea", label: "Nicaea II 787", date: 787, size: 6, extinct: false, uncertain: false,
                              event: true, eventKind: "Ecumenical council",
                              reference: wiki("Second_Council_of_Nicaea"),
                              doctrine: "The seventh ecumenical council, the last recognised by both East and West. It restored the veneration of icons, distinguishing honour from the worship due to God alone.",
                              children: [
                                { id: "east-west-schism", name: "East–West Schism", label: "Great Schism 1054", date: 1054, size: 9, extinct: false, uncertain: false,
                                  event: true, eventKind: "Schism",
                                  reference: wiki("East%E2%80%93West_Schism"),
                                  doctrine: "The long estrangement of the Greek East and the Latin West over papal authority, the filioque and differing customs. In 1054 legates and patriarch excommunicated each other, and the sack of Constantinople in 1204 sealed the split.",
                                  children: [
                                    { id: "latin-church", name: "Latin Church", date: 1054, size: 10, extinct: false, uncertain: false,
                                      tradition: "catholic",
                                      reference: wiki("Latin_Church"),
                                      doctrine: "The Western church under the bishop of Rome: a centralised papacy, scholastic theology and the mendicant orders.",
                                      children: [
                                        { id: "waldensians", name: "Waldensians", date: 1173, circa: true, size: 3, extinct: false, uncertain: false,
                                          tradition: "precursor",
                                          reference: wiki("Waldensians"),
                                          doctrine: "Followers of Peter Waldo of Lyon: poverty, lay preaching and scripture in the vernacular. Condemned in 1184, they survived in the Alps and adopted Reformed theology at Chanforan in 1532." },
                                        { id: "lollards", name: "Lollards", date: 1380, circa: true, end: 1530, endCirca: true, size: 3, extinct: true, uncertain: false,
                                          tradition: "precursor",
                                          reference: wiki("Lollardy"),
                                          doctrine: "English followers of John Wycliffe, who attacked clerical wealth and transubstantiation and championed an English Bible. Driven underground after 1414, they flowed into the English Reformation." },
                                        { id: "hussites", name: "Hussites", date: 1415, size: 4, extinct: false, uncertain: false,
                                          tradition: "precursor",
                                          reference: wiki("Hussites"),
                                          doctrine: "The Bohemian reform movement that arose after Jan Hus was burned at the Council of Constance in 1415. It demanded communion in both kinds and free preaching.",
                                          children: [
                                            { id: "taborites", name: "Taborites", date: 1420, end: 1434, size: 2, extinct: true, uncertain: false,
                                              reference: wiki("Taborites"),
                                              doctrine: "Radical, millenarian Hussites, militarily formidable under Jan Žižka. Utraquist forces destroyed them at the Battle of Lipany in 1434." },
                                            { id: "utraquists", name: "Utraquists", date: 1420, end: 1620, endCirca: true, size: 3, extinct: true, uncertain: false,
                                              reference: wiki("Utraquism"),
                                              doctrine: "Moderate Hussites who kept communion ‘in both kinds’ (sub utraque specie), recognised by the Compacts of Basel. They were absorbed after the Battle of White Mountain in 1620." },
                                            { id: "unity-of-brethren", name: "Unity of the Brethren", date: 1457, size: 3, extinct: false, uncertain: false,
                                              reference: wiki("Unity_of_the_Brethren"),
                                              doctrine: "The Unitas Fratrum, a peaceable, Bible-centred community that grew from the Hussites at Kunvald. Scattered after 1620, it was renewed at Herrnhut in 1722 as the Moravian Church." }
                                          ] },
                                        { id: "reformation", name: "Reformation", label: "Reformation 1517", date: 1517, size: 9, extinct: false, uncertain: false,
                                          event: true, eventKind: "Movement",
                                          reference: wiki("Reformation"),
                                          doctrine: "Luther’s Ninety-five Theses (1517) set off a break over justification by faith, the authority of scripture and the papacy that divided Western Christianity.",
                                          children: [
                                            { id: "nontrinitarians", name: "Nontrinitarians", date: 1550, circa: true, size: 3, extinct: false, uncertain: false,
                                              tradition: "radical",
                                              reference: wiki("Socinianism"),
                                              doctrine: "Radical reformers who rejected the Trinity as unscriptural, inspired by Michael Servetus (burned at Geneva in 1553) and by Lelio and Fausto Sozzini.",
                                              children: [
                                                { id: "polish-brethren", name: "Polish Brethren", date: 1565, end: 1658, size: 3, extinct: true, uncertain: false,
                                                  reference: wiki("Polish_Brethren"),
                                                  doctrine: "The Socinian church of Poland–Lithuania, centred on the academy at Raków and its Racovian Catechism. It was expelled from Poland in 1658." },
                                                { id: "transylvanian-unitarians", name: "Transylvanian Unitarians", date: 1568, size: 2, extinct: false, uncertain: false,
                                                  reference: wiki("Unitarian_Church_of_Transylvania"),
                                                  doctrine: "Founded by Ferenc Dávid after the Edict of Torda (1568), which granted religious toleration. It is the oldest Unitarian church still in existence." }
                                              ] },
                                            { id: "anabaptism", name: "Anabaptism", date: 1525, size: 5, extinct: false, uncertain: false,
                                              tradition: "anabaptist",
                                              reference: wiki("Anabaptism"),
                                              doctrine: "The ‘rebaptisers’ baptised only confessing believers and stood for a church free of the state, nonresistance and costly discipleship. The first adult baptisms were in Zurich in January 1525.",
                                              children: [
                                                { id: "swiss-brethren", name: "Swiss Brethren", date: 1525, size: 3, extinct: false, uncertain: false,
                                                  reference: wiki("Swiss_Brethren"),
                                                  doctrine: "The first Anabaptist congregation, gathered around Conrad Grebel, Felix Manz and George Blaurock. Its principles were set out in the Schleitheim Confession (1527).",
                                                  children: [
                                                    { id: "swiss-mennonites", name: "Swiss Mennonites", date: 1527, size: 3, extinct: false, uncertain: false,
                                                      reference: wiki("Swiss_Brethren"),
                                                      doctrine: "The continuing Swiss Brethren congregations, later known as Mennonites, persecuted in Bern and Zurich and scattered to the Jura, Alsace and the Palatinate." },
                                                    { id: "amish", name: "Amish", date: 1693, size: 3, extinct: false, uncertain: false,
                                                      reference: wiki("Amish"),
                                                      doctrine: "Followers of Jakob Ammann, who split from the Swiss Brethren and Alsatian Mennonites in 1693 over shunning and plainness of dress and life." }
                                                  ] },
                                                { id: "hutterites", name: "Hutterites", date: 1528, size: 3, extinct: false, uncertain: false,
                                                  reference: wiki("Hutterites"),
                                                  doctrine: "Communal Anabaptists of Moravia who held all goods in common, named after Jakob Hutter, who was martyred in 1536." },
                                                { id: "munster", name: "Münster Anabaptists", date: 1534, end: 1535, size: 2, extinct: true, uncertain: false,
                                                  reference: wiki("M%C3%BCnster_rebellion"),
                                                  doctrine: "Apocalyptic Anabaptists under Jan Matthys and John of Leiden seized Münster as a ‘New Jerusalem’. The city was retaken and the movement crushed in 1535." },
                                                { id: "mennonites", name: "Mennonites", date: 1536, size: 4, extinct: false, uncertain: false,
                                                  reference: wiki("Mennonites"),
                                                  doctrine: "Dutch and North German Anabaptists gathered by the former priest Menno Simons, who left Rome in 1536 and held them to peaceable nonresistance after Münster." }
                                              ] },
                                            { id: "anglicanism", name: "Anglicanism", date: 1534, size: 7, extinct: false, uncertain: false,
                                              tradition: "anglican",
                                              reference: wiki("Anglicanism"),
                                              doctrine: "The Act of Supremacy (1534) made the English monarch head of the church. Under Cranmer’s Book of Common Prayer it became a Reformed church that kept its catholic order.",
                                              children: [
                                                { id: "puritans", name: "Puritans", date: 1560, circa: true, size: 5, extinct: false, uncertain: false,
                                                  reference: wiki("Puritans"),
                                                  doctrine: "English Calvinists who sought to ‘purify’ the Church of England of remaining Catholic ceremony and of bishops. The movement broke apart into separate churches.",
                                                  children: [
                                                    { id: "quakers", name: "Quakers", date: 1652, size: 4, extinct: false, uncertain: false,
                                                      reference: wiki("Quakers"),
                                                      doctrine: "The Religious Society of Friends, gathered by George Fox in 1652. It teaches the inward light of Christ in every person and worships in silence, without clergy or sacraments." },
                                                    { id: "separatists", name: "English Separatists", date: 1581, size: 3, extinct: false, uncertain: false,
                                                      reference: wiki("Brownists"),
                                                      doctrine: "Puritans who despaired of reforming the national church and formed independent congregations. Robert Browne’s Norwich church of 1581 gave them the name Brownists.",
                                                      children: [
                                                        { id: "baptists", name: "Baptists", date: 1609, size: 5, extinct: false, uncertain: false,
                                                          reference: wiki("Baptists"),
                                                          doctrine: "John Smyth and Thomas Helwys, Separatist exiles in Amsterdam, founded the first Baptist church in 1609 on believer’s baptism and the freedom of each congregation.",
                                                          children: [
                                                            { id: "particular-baptists", name: "Particular Baptists", date: 1638, circa: true, size: 4, extinct: false, uncertain: false,
                                                              reference: wiki("Particular_Baptists"),
                                                              doctrine: "Calvinist Baptists, holding that Christ died for the elect in particular, who emerged from London’s semi-separatist congregations in the 1630s. They issued the First London Confession in 1644." },
                                                            { id: "general-baptists", name: "General Baptists", date: 1609, size: 3, extinct: false, uncertain: false,
                                                              reference: wiki("General_Baptists"),
                                                              doctrine: "Arminian Baptists in the line of Smyth and Helwys, holding that Christ died for all. Helwys founded the first Baptist church on English soil at Spitalfields in 1612." }
                                                          ] },
                                                        { id: "congregationalists", name: "Congregationalists", date: 1581, size: 4, extinct: false, uncertain: false,
                                                          reference: wiki("Congregational_church"),
                                                          doctrine: "Independent churches in which each congregation governs itself under Christ. Their heirs include the Pilgrims and the churches of Puritan New England." }
                                                      ] }
                                                  ] },
                                                { id: "church-of-england", name: "Church of England", date: 1559, size: 6, extinct: false, uncertain: false,
                                                  reference: wiki("Church_of_England"),
                                                  doctrine: "The national church as settled under Elizabeth I: the Acts of Supremacy and Uniformity, the Thirty-nine Articles and the Prayer Book, with bishops and a broadly Reformed theology." }
                                              ] },
                                            { id: "reformed", name: "Reformed", date: 1519, size: 8, extinct: false, uncertain: false,
                                              tradition: "reformed",
                                              reference: wiki("Reformed_Christianity"),
                                              doctrine: "The Swiss Reformation begun by Huldrych Zwingli in Zurich: the sovereignty of God, a spiritual view of the Lord’s Supper, and worship stripped to what scripture commands.",
                                              children: [
                                                { id: "zwinglians", name: "Zwinglians", date: 1523, end: 1549, size: 3, extinct: true, uncertain: false,
                                                  reference: wiki("Zwinglianism"),
                                                  doctrine: "Zwingli’s memorial view of the Eucharist, defended against Luther at Marburg in 1529. It merged with Calvin’s teaching in the Consensus Tigurinus of 1549." },
                                                { id: "calvinism", name: "Calvinism", date: 1541, size: 7, extinct: false, uncertain: false,
                                                  reference: wiki("Calvinism"),
                                                  doctrine: "John Calvin’s Geneva, to which he returned in 1541, became the model of Reformed theology and church order, spread through the Institutes and the Genevan Academy.",
                                                  children: [
                                                    { id: "remonstrants", name: "Remonstrants", date: 1610, size: 3, extinct: false, uncertain: false,
                                                      reference: wiki("Remonstrants"),
                                                      doctrine: "Dutch followers of Jacobus Arminius, who affirmed conditional election and resistible grace in the Remonstrance of 1610. Condemned at the Synod of Dort (1618–19), they survived as a separate church." },
                                                    { id: "dutch-reformed", name: "Dutch Reformed", date: 1571, size: 4, extinct: false, uncertain: false,
                                                      reference: wiki("Dutch_Reformed_Church"),
                                                      doctrine: "The Calvinist church of the Netherlands, organised at the Synod of Emden (1571) and bound to the Belgic Confession, the Heidelberg Catechism and the Canons of Dort." },
                                                    { id: "hungarian-reformed", name: "Hungarian Reformed", date: 1567, size: 3, extinct: false, uncertain: false,
                                                      reference: wiki("Reformed_Church_in_Hungary"),
                                                      doctrine: "The Calvinist church of Hungary and Transylvania, which adopted the Second Helvetic Confession at Debrecen in 1567." },
                                                    { id: "presbyterians", name: "Scottish Presbyterians", date: 1560, size: 5, extinct: false, uncertain: false,
                                                      reference: wiki("Presbyterianism"),
                                                      doctrine: "The Reformation Parliament of 1560 and John Knox’s Scots Confession created a church governed by elders in presbyteries rather than by bishops." },
                                                    { id: "huguenots", name: "Huguenots", date: 1559, size: 4, extinct: false, uncertain: false,
                                                      reference: wiki("Huguenots"),
                                                      doctrine: "French Reformed Protestants, who adopted the Gallican Confession in 1559. They lived through the Wars of Religion, won toleration by the Edict of Nantes (1598) and lost it in 1685." }
                                                  ] }
                                              ] },
                                            { id: "lutheranism", name: "Lutheranism", date: 1517, size: 8, extinct: false, uncertain: false,
                                              tradition: "lutheran",
                                              reference: wiki("Lutheranism"),
                                              doctrine: "Martin Luther’s reform: justification by grace through faith alone, scripture as supreme authority and Christ’s real presence in the Supper, summed up in the Augsburg Confession (1530).",
                                              children: [
                                                { id: "pietism", name: "Pietism", date: 1675, size: 4, extinct: false, uncertain: false,
                                                  reference: wiki("Pietism"),
                                                  doctrine: "A movement for heartfelt, practical faith, small-group Bible study and moral renewal, launched by Philipp Jakob Spener’s Pia Desideria in 1675." },
                                                { id: "church-of-sweden", name: "Church of Sweden", date: 1527, size: 4, extinct: false, uncertain: false,
                                                  reference: wiki("Church_of_Sweden"),
                                                  doctrine: "Gustav Vasa broke with Rome at the Diet of Västerås in 1527. The church kept its bishops and confirmed the Augsburg Confession at the Uppsala Synod in 1593." },
                                                { id: "confessional-lutherans", name: "Confessional Lutherans", date: 1580, size: 6, extinct: false, uncertain: false,
                                                  reference: wiki("Book_of_Concord"),
                                                  doctrine: "The main Lutheran churches of Germany and Scandinavia, united around the Book of Concord (1580), which collects the Lutheran confessions." },
                                                { id: "philippists", name: "Philippists", date: 1548, end: 1580, size: 2, extinct: true, uncertain: false,
                                                  reference: wiki("Philippists"),
                                                  doctrine: "Followers of Philip Melanchthon, who accepted the Leipzig Interim and moved toward Reformed views of the Supper. The Formula of Concord (1577) sidelined them and some became Reformed." },
                                                { id: "gnesio-lutherans", name: "Gnesio-Lutherans", date: 1548, end: 1580, size: 2, extinct: true, uncertain: false,
                                                  reference: wiki("Gnesio-Lutherans"),
                                                  doctrine: "The ‘genuine Lutherans’ led by Matthias Flacius, who resisted the Interims and any softening of Luther’s teaching. Their concerns were absorbed into the Formula of Concord." }
                                              ] },
                                            { id: "roman-catholic", name: "Roman Catholic Church", date: 1545, size: 10, extinct: false, uncertain: false,
                                              tradition: "catholic",
                                              reference: wiki("Catholic_Church"),
                                              doctrine: "The Western church as renewed by the Council of Trent (1545–63). It reaffirmed the seven sacraments, transubstantiation, and scripture with tradition, and launched the Catholic Reformation with the Jesuits." }
                                          ] }
                                      ] },
                                    { id: "eastern-orthodox", name: "Eastern Orthodox Church", date: 1054, size: 9, extinct: false, uncertain: false,
                                      tradition: "orthodox",
                                      reference: wiki("Eastern_Orthodox_Church"),
                                      doctrine: "The Greek and Slavic churches in communion with the Ecumenical Patriarch. They hold to the seven councils, conciliar authority, the creed without the filioque, and the Byzantine liturgy.",
                                      children: [
                                        { id: "constantinople-patriarchate", name: "Patriarchate of Constantinople", label: "Ecumenical Patriarchate", date: 1054, size: 7, extinct: false, uncertain: false,
                                          reference: wiki("Ecumenical_Patriarchate_of_Constantinople"),
                                          doctrine: "The Greek Orthodox Church centred on Constantinople, whose patriarch is ‘first among equals’ among Orthodox bishops. It continued under Ottoman rule after 1453." },
                                        { id: "serbian-orthodox", name: "Serbian Orthodox", date: 1219, size: 4, extinct: false, uncertain: false,
                                          reference: wiki("Serbian_Orthodox_Church"),
                                          doctrine: "Saint Sava won self-governance (autocephaly) for the Serbian church in 1219. It became a patriarchate in 1346." },
                                        { id: "union-of-florence", name: "Union of Florence", date: 1439, end: 1484, size: 2, extinct: true, uncertain: false,
                                          reference: wiki("Council_of_Florence"),
                                          doctrine: "A reunion of East and West agreed at Florence in 1439 under the shadow of the Ottomans. Most Orthodox clergy and laity rejected it, and a synod at Constantinople repudiated it in 1484." },
                                        { id: "russian-orthodox", name: "Russian Orthodox", date: 1448, size: 6, extinct: false, uncertain: false,
                                          reference: wiki("Russian_Orthodox_Church"),
                                          doctrine: "The Russian church elected its own metropolitan in 1448 after rejecting the Union of Florence. Moscow became a patriarchate in 1589 and saw itself as the ‘Third Rome’." },
                                        { id: "union-of-brest", name: "Ruthenian Greek Catholics", date: 1596, size: 3, extinct: false, uncertain: false,
                                          reference: wiki("Union_of_Brest"),
                                          doctrine: "At the Union of Brest (1596), Orthodox bishops of Poland–Lithuania entered communion with Rome while keeping the Byzantine rite. Their heir is today’s Ukrainian Greek Catholic Church." }
                                      ] }
                                  ] }
                              ] },
                            { id: "iconoclasm", name: "Byzantine Iconoclasm", label: "Iconoclasm", date: 726, circa: true, end: 843, size: 4, extinct: true, uncertain: false,
                              reference: wiki("Byzantine_Iconoclasm"),
                              doctrine: "The imperial destruction of religious images begun under Leo III. Reversed at Nicaea in 787, it returned in 815 and ended with the ‘Triumph of Orthodoxy’ in 843." },
                            { id: "paulicians", name: "Paulicians", date: 650, circa: true, end: 900, endCirca: true, size: 3, extinct: true, uncertain: false,
                              reference: wiki("Paulicianism"),
                              doctrine: "A dualist sect of Armenia and Asia Minor that rejected the Byzantine sacraments, clergy and icons. Its stronghold at Tephrike fell to Byzantium in 872.",
                              children: [
                                { id: "bogomils", name: "Bogomils", date: 940, circa: true, end: 1450, endCirca: true, size: 3, extinct: true, uncertain: true,
                                  reference: wiki("Bogomilism"),
                                  doctrine: "A Balkan dualist movement begun by the priest Bogomil in Bulgaria, rejecting the material world, the sacraments and the hierarchy. How far it descends from the Paulicians is debated.",
                                  children: [
                                    { id: "cathars", name: "Cathars", date: 1143, circa: true, end: 1321, size: 4, extinct: true, uncertain: true,
                                      reference: wiki("Catharism"),
                                      doctrine: "Dualist ‘good Christians’ of Languedoc and northern Italy, led by ascetic perfecti. The Albigensian Crusade and the Inquisition destroyed them; their debt to the Bogomils is debated." }
                                  ] }
                              ] },
                            { id: "maronites", name: "Maronite Church", date: 685, circa: true, size: 3, extinct: false, uncertain: false,
                              tradition: "catholic",
                              reference: wiki("Maronite_Church"),
                              doctrine: "The Syriac Chalcedonian church of Lebanon, named for Saint Maron. It elected its own patriarch c. 685 and confirmed full communion with Rome in 1182." },
                            { id: "monothelitism", name: "Monothelitism", date: 633, end: 712, endCirca: true, size: 3, extinct: true, uncertain: false,
                              reference: wiki("Monothelitism"),
                              doctrine: "A compromise teaching that Christ has two natures but a single will, promoted by Emperor Heraclius. The sixth ecumenical council (Constantinople III, 681) condemned it." }
                          ] },
                        { id: "oriental-orthodoxy", name: "Oriental Orthodoxy", date: 451, size: 7, extinct: false, uncertain: false,
                          tradition: "oriental",
                          reference: wiki("Oriental_Orthodox_Churches"),
                          doctrine: "The churches that rejected Chalcedon’s ‘two natures’. Following Cyril of Alexandria they confess one united divine-human nature of the incarnate Word (miaphysitism), and they too condemn Eutyches.",
                          children: [
                            { id: "coptic-orthodox", name: "Coptic Orthodox", date: 451, size: 5, extinct: false, uncertain: false,
                              reference: wiki("Coptic_Orthodox_Church"),
                              doctrine: "The Church of Alexandria that stayed loyal to the deposed Patriarch Dioscorus: the heart of Egyptian Christianity and the cradle of monasticism." },
                            { id: "ethiopian-orthodox", name: "Ethiopian Tewahedo", date: 480, circa: true, size: 4, extinct: false, uncertain: false,
                              reference: wiki("Ethiopian_Orthodox_Tewahedo_Church"),
                              doctrine: "Christian since King Ezana of Aksum (c. 330), the church took its miaphysite character through Alexandria and the Nine Saints of the late fifth century. Tewahedo means ‘made one’." },
                            { id: "armenian-apostolic", name: "Armenian Apostolic", date: 506, size: 4, extinct: false, uncertain: false,
                              reference: wiki("Armenian_Apostolic_Church"),
                              doctrine: "Armenia adopted Christianity as a state c. 301. Its church formally rejected Chalcedon at the synods of Dvin in 506 and 554." },
                            { id: "julianists", name: "Julianists", date: 520, circa: true, end: 700, endCirca: true, size: 2, extinct: true, uncertain: false,
                              reference: wiki("Aphthartodocetae"),
                              doctrine: "Followers of Julian of Halicarnassus, who held that Christ’s body was incorruptible from the moment of the incarnation. The dispute divided the miaphysites against Severus of Antioch." },
                            { id: "syriac-orthodox", name: "Syriac Orthodox", date: 543, circa: true, size: 4, extinct: false, uncertain: false,
                              reference: wiki("Syriac_Orthodox_Church"),
                              doctrine: "The miaphysite church of Syria, whose separate hierarchy was built by Jacob Baradaeus’ secret ordinations from c. 543, hence the old name ‘Jacobite’." }
                          ] },
                        { id: "eutychianism", name: "Eutychianism", date: 451, end: 500, endCirca: true, size: 3, extinct: true, uncertain: false,
                          reference: wiki("Eutychianism"),
                          doctrine: "Eutyches taught that Christ’s humanity was absorbed into his divinity, leaving a single nature. Condemned at Chalcedon, and also by the Oriental Orthodox churches." }
                      ] },
                    { id: "nestorianism", name: "Nestorianism", date: 431, size: 6, extinct: false, uncertain: false,
                      tradition: "east",
                      reference: wiki("Nestorianism"),
                      doctrine: "The Antiochene stress on Christ’s two distinct natures associated with Nestorius, who preferred ‘Christ-bearer’ to ‘God-bearer’ for Mary. The Church of the East rejects the ‘Nestorian’ label.",
                      children: [
                        { id: "church-of-the-east", name: "Church of the East", date: 486, size: 6, extinct: false, uncertain: false,
                          reference: wiki("Church_of_the_East"),
                          doctrine: "The church of the Persian Empire, which adopted a two-nature Christology at its synods of 484–486 and carried the faith along the Silk Road to Central Asia, China and India.",
                          children: [
                            { id: "jingjiao", name: "Jingjiao (China)", date: 635, end: 1368, endCirca: true, size: 3, extinct: true, uncertain: false,
                              reference: wiki("Church_of_the_East_in_China"),
                              doctrine: "The Church of the East in China: the monk Alopen reached Chang’an in 635 and the Xi’an Stele was raised in 781. It faded after the Tang, revived under the Mongols and vanished with the Yuan." },
                            { id: "saint-thomas-christians", name: "Saint Thomas Christians", date: 650, circa: true, size: 3, extinct: false, uncertain: false,
                              reference: wiki("Saint_Thomas_Christians"),
                              doctrine: "The ancient Christians of Kerala, who trace their faith to the apostle Thomas and were long served by bishops of the Church of the East. They divided after 1653, and some joined the Syriac Orthodox." },
                            { id: "chaldean-catholic", name: "Chaldean Catholic Church", date: 1552, size: 4, extinct: false, uncertain: false,
                              reference: wiki("Chaldean_Catholic_Church"),
                              doctrine: "Formed when part of the Church of the East, under Patriarch Shimun VIII Yohannan Sulaqa, entered communion with Rome in 1553." },
                            { id: "assyrian-church", name: "Assyrian Church of the East", date: 1552, size: 4, extinct: false, uncertain: false,
                              reference: wiki("Assyrian_Church_of_the_East"),
                              doctrine: "The continuing Church of the East outside communion with Rome. After the schism of 1552 its patriarchal lines moved back and forth between independence and union." }
                          ] }
                      ] }
                  ] }
              ] }
          ] }
      ] }
  ]
};

/* =====================================================================
   2. CONFIG
   ===================================================================== */
const TAU = Math.PI * 2;
const R = 420;                          // radius of the 1700 ring (the bark)
const MIN_GAP = 12;                     // a child always sits this far beyond its parent
const TIME = d3.scaleLinear()           // year -> radius, piecewise to give eventful eras room
  .domain([49, 150, 325, 451, 700, 1054, 1300, 1517, 1700])
  .range([0, .14, .28, .37, .46, .565, .665, .78, 1].map(f => f * R))
  .clamp(true);
const LEAF_BASE = R + 3;                // living lineages end in a leaf just outside the bark
const LEAF_LABEL_R = R + 28;
const SAMPLES = 28;                     // points per limb
const RING_LABEL_ANGLE = 214 * Math.PI / 180;
const INSC_FS = 8.5;                    // lettering in the gilt band
const INSC_HALF = 9;                    // half the width of that band
const ROMAN = { 100: "C", 200: "CC", 300: "CCC", 400: "CD", 500: "D", 700: "DCC", 1000: "M", 1300: "MCCC", 1500: "MD" };
const REDUCED = matchMedia("(prefers-reduced-motion: reduce)").matches;
const DUR = REDUCED ? 0 : 850, STAGGER = REDUCED ? 0 : 55;
const INTRO_DUR = REDUCED ? 0 : 1250, INTRO_STAGGER = REDUCED ? 0 : 45;
const nodeR = d3.scaleSqrt().domain([1, 10]).range([2, 7]);

const TRADITIONS = {
  trunk:      { family: "trunk",      label: "Great Church (undivided)", shift: 0 },
  east:       { family: "east",       label: "Church of the East",       shift: 0 },
  oriental:   { family: "oriental",   label: "Oriental Orthodox",        shift: 0 },
  orthodox:   { family: "orthodox",   label: "Eastern Orthodox",         shift: 0 },
  catholic:   { family: "catholic",   label: "Catholic",                 shift: 0 },
  precursor:  { family: "protestant", label: "Pre-Reformation reform",   shift: 7 },
  lutheran:   { family: "protestant", label: "Lutheran",                 shift: 0 },
  reformed:   { family: "protestant", label: "Reformed",                 shift: -8 },
  anglican:   { family: "protestant", label: "Anglican",                 shift: 4 },
  anabaptist: { family: "protestant", label: "Anabaptist",               shift: -4 },
  radical:    { family: "protestant", label: "Radical Reformation",      shift: 8 }
};
const FAMILIES = {
  trunk: "Great Church", east: "Church of the East", oriental: "Oriental Orthodox",
  orthodox: "Eastern Orthodox", catholic: "Catholic", protestant: "Protestant"
};

let PAL = {};
function readPalette() {
  const cs = getComputedStyle(document.documentElement);
  const v = name => cs.getPropertyValue(name).trim();
  PAL = { bg: v("--bg"), extinct: v("--extinct"), gold0: v("--gold-0"), gold1: v("--gold-1"), trad: {} };
  for (const f of Object.keys(FAMILIES)) PAL[f] = v("--" + f);
  for (const [key, t] of Object.entries(TRADITIONS)) {
    const c = d3.hcl(PAL[t.family]);
    c.l += t.shift;
    PAL.trad[key] = c.formatHex();
  }
}

// Seeded randomness, so the hand-drawn irregularities are the same on every visit.
function mulberry32(seed) {
  return function () {
    seed = seed + 0x6D2B79F5 | 0;
    let t = Math.imul(seed ^ seed >>> 15, 1 | seed);
    t = t + Math.imul(t ^ t >>> 7, 61 | t) ^ t;
    return ((t ^ t >>> 14) >>> 0) / 4294967296;
  };
}
function hashSeed(s) {
  let h = 2166136261;
  for (const ch of s) h = Math.imul(h ^ ch.charCodeAt(0), 16777619);
  return h >>> 0;
}

/* =====================================================================
   3. VALIDATE
   ===================================================================== */
function validate(root) {
  const warn = (n, msg) => console.warn(`[cladogram] ${n.data.id}: ${msg}`);
  const ids = new Set();
  root.each(n => {
    const d = n.data;
    for (const f of ["id", "name", "date", "doctrine", "size", "extinct", "uncertain", "reference"]) {
      if (!(f in d)) warn(n, `missing field "${f}"`);
    }
    if (ids.has(d.id)) warn(n, "duplicate id");
    ids.add(d.id);
    if (n.parent && d.date < n.parent.data.date) warn(n, `date ${d.date} precedes parent's ${n.parent.data.date}`);
    if (d.extinct && !(d.end >= d.date)) warn(n, "extinct node needs end ≥ date");
    if (n.parent && n.parent.data.extinct && !d.extinct) warn(n, "descendant of an extinct node must be extinct");
    if (!d.extinct && n.children && !n.children.some(c => !c.data.extinct)) warn(n, "surviving internal node has no surviving child");
    if (!(d.size >= 1 && d.size <= 10)) warn(n, "size outside 1–10");
    if (d.reference !== null && !/^https:\/\/en\.wikipedia\.org\/wiki\/\S+$/.test(d.reference)) warn(n, "reference is not a Wikipedia URL");
    if (d.tradition && !TRADITIONS[d.tradition]) warn(n, `unknown tradition "${d.tradition}"`);
  });
  console.info(`[cladogram] ${ids.size} nodes, ${root.leaves().length} leaves`);
}

/* =====================================================================
   4. PREP
   ===================================================================== */
const root = d3.hierarchy(DATA);
validate(root);

const ALL = [];
root.eachBefore(n => ALL.push(n));      // pre-order: DOM order and tab order follow the tree

root.eachBefore(n => {
  n.allKids = n.children || null;
  n.collapsed = false;
  n.trad = n.data.tradition || (n.parent ? n.parent.trad : "trunk");
  n.family = TRADITIONS[n.trad].family;
  n.nr = nodeR(n.data.size);
  n.r = n.parent ? Math.max(TIME(n.data.date), n.parent.r + MIN_GAP) : 0;
  n.dimT = n.dimC = 1;
  const rnd = mulberry32(hashSeed(n.data.id));
  n.easeP = 1.9 + rnd() * 0.7;          // each limb bends a little differently,
  n.sway = (rnd() - 0.5) * 6;           // bows slightly to one side (px),
  n.tilt = (rnd() - 0.5) * 26;          // and no leaf sits perfectly straight (deg)
});
root.eachAfter(n => {
  const kids = n.allKids || [];
  n.mass = n.data.size + d3.sum(kids, c => c.mass);
  n.leafTotal = kids.length ? d3.sum(kids, c => c.leafTotal) : 1;
  n.descCount = d3.sum(kids, c => c.descCount + 1);
  n.lastEnd = n.data.extinct ? Math.max(n.data.end, d3.max(kids, c => c.lastEnd) || -Infinity) : Infinity;
});
const width = d3.scaleSqrt().domain([1, d3.max(root.children, c => c.mass)]).range([0.7, 4.2]);

/* =====================================================================
   5. LAYOUT  (angle from a weighted partition, radius from time)
   ===================================================================== */
const partition = d3.partition().size([TAU, 1]);
let rotation = null;
const terminalR = n => n.data.extinct ? TIME(n.lastEnd) : R;

function leafWeight(n) {
  // leaves whose labels sit near the centre get a wider wedge, evening out arc length
  const f = Math.pow(R / Math.max(terminalR(n), 0.22 * R), 0.9);
  return (n.allKids ? 1 + 0.5 * Math.sqrt(n.leafTotal) : 1) * f;
}

function layout() {
  for (const n of ALL) {
    n.children = n.allKids && !n.collapsed ? n.allKids : undefined;
    n.visible = !n.parent || (n.parent.visible && !n.parent.collapsed);
  }
  root.eachAfter(n => { n.value = n.children ? d3.sum(n.children, c => c.value) : leafWeight(n); });
  partition(root);
  // Turn the tree so its living canopy (surviving lineages, weighted by wedge) faces north and the
  // short-lived early branches settle below like roots; Latin West falls left, the East right.
  // Fixed from the full tree so folding a branch never swings the whole picture round.
  if (rotation === null) {
    let sx = 0, sy = 0;
    root.leaves().forEach(n => {
      if (n.data.extinct) return;
      const t = (n.x0 + n.x1) / 2, w = n.x1 - n.x0;
      sx += Math.sin(t) * w;
      sy += Math.cos(t) * w;
    });
    rotation = -Math.atan2(sx, sy);
  }
  for (const n of ALL) {
    if (n.visible) { n.tt = (n.x0 + n.x1) / 2 + rotation; n.rt = n.r; }
    else { n.tt = n.parent.tt; n.rt = n.parent.rt; }   // hidden nodes rest inside their visible ancestor
  }
}

/* =====================================================================
   6. GEOMETRY
   ===================================================================== */
const pt = (t, r) => [r * Math.sin(t), -r * Math.cos(t)];
const norm = t => ((t % TAU) + TAU) % TAU;
const centerLine = d3.line().curve(d3.curveBasis);
const sideLine = d3.line().curve(d3.curveCatmullRom.alpha(0.5));

// A limb peels off its parent sideways, then runs outward; `sway` bows it gently to one side.
function limbPoints(t0, r0, t1, r1, p, sway) {
  const pts = new Array(SAMPLES + 1);
  const span = Math.min(1, (Math.abs(r1 - r0) + Math.abs(t1 - t0) * (r0 + r1) / 2) / 90);
  for (let i = 0; i <= SAMPLES; i++) {
    const u = i / SAMPLES;
    const t = t0 + (t1 - t0) * (1 - Math.pow(1 - u, p));
    const r = r0 + (r1 - r0) * u;
    const b = sway * span * Math.sin(Math.PI * u);
    pts[i] = [r * Math.sin(t) + b * Math.cos(t), -r * Math.cos(t) + b * Math.sin(t)];
  }
  return pts;
}

function limbOf(n) {
  const c = n.cur, p = n.parent.cur;
  const t0 = n.parent.parent ? p.t : c.t;               // limbs leave the root radially
  return limbPoints(t0, p.r, c.t, c.r, n.easeP, n.sway);
}

// A filled, tapering ribbon: full width at the base, narrowing quickly toward the tip.
function ribbonPath(pts, w0, w1) {
  const left = [], right = [], last = pts.length - 1;
  for (let i = 0; i <= last; i++) {
    const a = pts[Math.max(0, i - 1)], b = pts[Math.min(last, i + 1)];
    const dx = b[0] - a[0], dy = b[1] - a[1], len = Math.hypot(dx, dy) || 1;
    const nx = -dy / len, ny = dx / len;
    const w = (w1 + (w0 - w1) * Math.pow(1 - i / last, 1.4)) / 2;
    left.push([pts[i][0] + nx * w, pts[i][1] + ny * w]);
    right.push([pts[i][0] - nx * w, pts[i][1] - ny * w]);
  }
  return sideLine(left) + "L" + sideLine(right.reverse()).slice(1) + "Z";
}

function arcPath(t0, t1, r, reverse) {
  const [xa, ya] = pt(reverse ? t1 : t0, r), [xb, yb] = pt(reverse ? t0 : t1, r);
  return `M${xa},${ya}A${r},${r} 0 ${t1 - t0 > Math.PI ? 1 : 0} ${reverse ? 0 : 1} ${xb},${yb}`;
}

/* =====================================================================
   7. RENDER
   ===================================================================== */
const svg = d3.select("#chart");
const viewport = d3.select("#viewport");
const defs = d3.select("#defs");
const gWood = d3.select("#wood"), gInscBase = d3.select("#insc-base"), gInscFam = d3.select("#insc-fam");
const gWash = d3.select("#wash"), gTails = d3.select("#tails"), gLimbs = d3.select("#limbs");
const gLeaves = d3.select("#leaves"), gNodes = d3.select("#nodes"), gLabels = d3.select("#labels");

// A three-lobed ivy leaf, as in the bar-and-ivy borders of Gothic books; it points outward.
const IVY = "M0,-.4C-1.8,.5 -4.5,-.5 -5.9,-4.6C-4.5,-4.8 -3.3,-5.4 -2.7,-6.6C-2.4,-8.6 -1.2,-10.3 0,-12.3" +
  "C1.2,-10.3 2.4,-8.6 2.7,-6.6C3.3,-5.4 4.5,-4.8 5.9,-4.6C4.5,-.5 1.8,.5 0,-.4Z";
const IVY_VEIN = "M0,-1.2L0,-10.2M0,-3.6L-3.9,-4.7M0,-3.6L3.9,-4.7";
// a cross pattée for the root: four arms flaring from a narrow centre
function crossPath(s) {
  const n = s * 0.14, a = s * 0.42;
  return `M${-n},${-n}L${-a},${-s}L${a},${-s}L${n},${-n}L${s},${-a}L${s},${a}L${n},${n}L${a},${s}` +
    `L${-a},${s}L${-n},${n}L${-s},${a}L${-s},${-a}Z`;
}

const fmtYear = (y, circa) => (circa ? "c. " : "") + y;
function whenText(d) {
  if (d.extinct) return `${fmtYear(d.date, d.circa)} – ${fmtYear(d.end, d.endCirca)} · Extinct`;
  if (d.event) return `${fmtYear(d.date, d.circa)} · ${d.eventKind}`;
  return `${fmtYear(d.date, d.circa)} · Lineage continues`;
}
const labelText = n => {
  const s = n.data.label || n.data.name;
  return n.data.event ? s.toUpperCase() : s;
};

for (const n of ALL) {
  const d = n.data;
  if (n.parent) n.limbEl = gLimbs.append("path").attr("class", "limb").node();
  if (n.parent && !d.extinct && !d.uncertain) n.washEl = gWash.append("path").node();
  n.tailEl = gTails.append("path").attr("class", "tail" + (d.extinct ? " ext" : "")).node();
  if (d.extinct) {
    const g = defs.append("linearGradient").attr("id", "fade-" + d.id).attr("gradientUnits", "userSpaceOnUse");
    g.append("stop").attr("offset", 0).style("stop-color", "var(--extinct)").style("stop-opacity", .8);
    g.append("stop").attr("offset", 1).style("stop-color", "var(--extinct)").style("stop-opacity", 0);
    n.gradEl = g.node();
  } else if (n.parent) {
    const lf = gLeaves.append("g").attr("class", "leaf");
    lf.append("path").attr("class", "blade").attr("d", IVY);
    lf.append("path").attr("class", "vein").attr("d", IVY_VEIN);
    n.leafEl = lf.node();
    n.leafScale = 0.55 + 0.24 * Math.sqrt(d.size);
  }
  const ng = gNodes.append("g").datum(n)
    .attr("class", "node" + (d.extinct ? " ext" : "") + (d.uncertain ? " unc" : "") + (d.event ? " event" : ""))
    .attr("tabindex", 0).attr("role", "button")
    .attr("aria-label", `${d.name}, ${whenText(d)}`);
  ng.append("circle").attr("class", "hit").attr("r", Math.max(n.nr + 5, 9));
  if (!d.extinct && !d.event) ng.append("circle").attr("class", "halo").attr("r", n.nr + 2.6);
  ng.append("circle").attr("class", "bud").attr("r", n.nr + 3.5);
  ng.append("circle").attr("class", "dot").attr("r", n.nr);
  if (!d.extinct && !d.event) ng.append("circle").attr("class", "gleam")       // a touch of lead white
    .attr("r", n.nr * 0.3).attr("cx", -n.nr * 0.34).attr("cy", -n.nr * 0.36);
  if (d.event && n.parent) ng.append("circle").attr("class", "knot").attr("r", n.nr * 0.45);
  if (!n.parent) ng.append("path").attr("class", "cross").attr("d", crossPath(n.nr * 0.66));
  if (d.uncertain) ng.append("text").attr("class", "q").attr("x", n.nr + 0.8).attr("y", -n.nr + 1.5).text("?");
  n.nodeEl = ng.node();
  const lw = gLabels.append("g").attr("class", "lw").datum(n);
  n.labelText = lw.append("text").text(labelText(n)).node();
  n.labelWrap = lw.node();
}

function colorOf(n) {
  return n.data.extinct ? PAL.extinct : PAL.trad[n.trad];
}

function paint() {
  readPalette();
  for (const n of ALL) {
    const d = n.data, c = colorOf(n);
    const soft = d.uncertain ? d3.interpolateRgb(c, PAL.bg)(0.4) : c;
    const edge = d3.hcl(c).darker(1.1).formatHex();      // the contour a limner inks round each painted shape
    if (n.limbEl) {
      n.ribbon = !d.extinct && !d.uncertain;
      d3.select(n.limbEl)
        .classed("ext", d.extinct && !d.uncertain).classed("unc", d.uncertain)
        .attr("fill", n.ribbon ? c : "none")
        .attr("stroke", n.ribbon ? edge : soft)
        .attr("stroke-width", n.ribbon ? 0.4 : null)
        .attr("stroke-opacity", n.ribbon ? 0.7 : null);
    }
    if (n.washEl) n.washEl.setAttribute("fill", c);
    if (n.leafEl) {
      const blade = n.leafEl.querySelector(".blade");
      blade.setAttribute("fill", c);
      blade.setAttribute("stroke", edge);
    }
    d3.select(n.tailEl)
      .attr("stroke", d.extinct ? `url(#fade-${d.id})` : c)
      .attr("stroke-width", d.extinct ? null : Math.max(0.8, width(n.mass) * 0.62));
    const ng = d3.select(n.nodeEl);
    if (d.event) ng.select(".dot").attr("fill", "url(#gilt)").attr("stroke", PAL.gold0);
    else if (d.extinct) ng.select(".dot").attr("fill", d3.interpolateRgb(PAL.extinct, PAL.bg)(0.3)).attr("stroke", PAL.extinct);
    else ng.select(".dot").attr("fill", soft).attr("stroke", d.uncertain ? c : edge);
    ng.select(".halo").attr("fill", c);
    ng.select(".knot").attr("stroke", PAL.gold0);
    ng.select(".bud").attr("stroke", d.event ? PAL.gold0 : c);
  }
}

/* ---------- the rota: glow, growth rings, the ruled edge and the gilt band ---------- */
const woodRnd = mulberry32(1054);
const HARMONICS = d3.range(5).map(k => ({ k: k + 2, a: (0.45 + woodRnd() * 0.55) / (k + 1), p: woodRnd() * TAU }));
const ROUGH = d3.range(3).map(k => ({ k: 29 + k * 17, a: 0.35 / (k + 1), p: woodRnd() * TAU }));
const waves = (list, t, salt) => list.reduce((s, h) => s + h.a * Math.sin(h.k * t + h.p + salt), 0);
const woodLine = d3.line().curve(d3.curveCardinalClosed);

// Rings share one gently irregular shape, as real growth rings do, each with its own small drift.
function ringPath(r0, salt, rough = 0) {
  const amp = 0.4 + 1.6 * r0 / R, count = Math.max(48, Math.round(r0 / 2.5));
  return woodLine(d3.range(count).map(i => {
    const t = i / count * TAU;
    return pt(t, r0 + amp * waves(HARMONICS, t, salt) + rough * waves(ROUGH, t, salt * 3));
  }));
}

const circlePath = r => `M0,${-r}A${r},${r} 0 1 1 0,${r}A${r},${r} 0 1 1 0,${-r}Z`;
const goldStops = (g, stops) => stops.forEach(([o, k]) => g.append("stop").attr("offset", o).style("stop-color", `var(--gold-${k})`));
let glintAnim = null;

// Gold leaf: a burnished sheen for each small piece, a wider one with several highlights for the band,
// and a single glint that passes over the band once the tree has grown.
function defineGilt() {
  goldStops(defs.append("linearGradient").attr("id", "gilt").attr("x1", 0).attr("y1", 0).attr("x2", 1).attr("y2", 1),
    [[0, 0], [.3, 1], [.48, 2], [.64, 3], [1, 4]]);
  goldStops(defs.append("linearGradient").attr("id", "gilt-wide").attr("x1", 0).attr("y1", 0).attr("x2", 1).attr("y2", 1),
    [[0, 0], [.14, 2], [.27, 1], [.4, 3], [.5, 2], [.62, 1], [.76, 3], [.87, 2], [1, 4]]);
  const gl = defs.append("linearGradient").attr("id", "glint").attr("gradientUnits", "userSpaceOnUse")
    .attr("x1", -90).attr("y1", -90).attr("x2", 90).attr("y2", 90).attr("gradientTransform", "translate(-1000,-1000)");
  for (const [o, a] of [[0, 0], [.5, .95], [1, 0]]) gl.append("stop").attr("offset", o).attr("stop-color", "#fff6dc").attr("stop-opacity", a);
  glintAnim = gl.append("animateTransform").attr("attributeName", "gradientTransform").attr("type", "translate")
    .attr("from", "-1000 -1000").attr("to", "1000 1000").attr("dur", "2.6s").attr("begin", "indefinite").attr("fill", "freeze").node();
}

function drawWood() {
  defineGilt();
  const glow = defs.append("radialGradient").attr("id", "pithGlow");
  glow.append("stop").attr("offset", 0).style("stop-color", "var(--glow)").style("stop-opacity", .9);
  glow.append("stop").attr("offset", .35).style("stop-color", "var(--glow)").style("stop-opacity", .45);
  glow.append("stop").attr("offset", 1).style("stop-color", "var(--glow)").style("stop-opacity", 0);
  gWood.append("circle").attr("class", "pith").attr("r", 0.34 * R).attr("fill", "url(#pithGlow)").attr("opacity", .75);

  const centuries = d3.range(100, 1701, 100);
  const salt = new Map(centuries.map(y => [y, (woodRnd() - 0.5) * 0.4]));
  // alternate centuries carry a faint tone, like early- and latewood
  for (let y = 100; y < 1700; y += 200) {
    gWood.append("path").attr("class", "wood-band").style("--d", y / 2 + "ms")
      .attr("d", ringPath(TIME(y + 100), salt.get(y + 100)) + ringPath(TIME(y), salt.get(y)));
  }
  // fine annual lines between the centuries
  for (let y = 60; y < 1700; y += 20) {
    if (y % 100 === 0) continue;
    gWood.append("path").attr("class", "wood-ring fine").attr("pathLength", 1).style("--d", y * 0.9 + "ms")
      .attr("d", ringPath(TIME(y), (woodRnd() - 0.5) * 0.6));
  }
  centuries.slice(0, -1).forEach(y => {
    gWood.append("path").attr("class", "wood-ring").attr("pathLength", 1).style("--d", y * 0.9 + "ms")
      .attr("d", ringPath(TIME(y), salt.get(y)));
  });
  // the edge at 1700, ruled twice with the compass: once in ink, once in red
  gWood.append("path").attr("class", "rule ink").attr("pathLength", 1).style("--d", "1500ms")
    .attr("d", ringPath(R, salt.get(1700), 0.3));
  gWood.append("path").attr("class", "rule red").attr("pathLength", 1).style("--d", "1650ms")
    .attr("d", ringPath(R + 5, salt.get(1700) + 0.2));

  // the gilt band that carries the inscription, edged outside with red and blue pearls
  const rim = gWood.append("g").attr("class", "rim").style("--d", "1500ms");
  const r0 = INSC_R - INSC_HALF, r1 = INSC_R + INSC_HALF;
  rim.append("path").attr("class", "insc-band").attr("d", circlePath(r1) + circlePath(r0));
  for (const r of [r0, r1]) {
    rim.append("circle").attr("class", "fillet").attr("r", r);
    rim.append("circle").attr("class", "fillet-edge").attr("r", r - 1.15);
    rim.append("circle").attr("class", "fillet-edge").attr("r", r + 1.15);
  }
  for (const r of [r0, r1]) rim.append("circle").attr("class", "glint").attr("r", r);
  const rp = r1 + 6, beads = Math.round(TAU * rp / 14);
  for (const [cls, off] of [["r", 0], ["b", 7]]) {
    rim.append("path").attr("class", "pearls " + cls).attr("d", circlePath(rp))
      .attr("pathLength", beads * 14).attr("stroke-dasharray", "0 14").attr("stroke-dashoffset", -off);
  }

  // centuries numbered in red, as a scribe would
  gWood.selectAll("text.ring-lbl").data([100, 200, 300, 400, 500, 700, 1000, 1300, 1500]).join("text")
    .attr("class", "ring-lbl").attr("text-anchor", "middle").attr("dy", "0.35em")
    .style("--d", y => y + 400 + "ms")
    .text(y => ROMAN[y])
    .attr("transform", ringLabelTransform);
}

function ringLabelTransform(y) {
  const [x, yy] = pt(RING_LABEL_ANGLE, TIME(y));
  return `translate(${x},${yy}) scale(${labelScale})`;
}

/* ---------- per-frame drawing from each node's current polar state ---------- */
let labelScale = 1, currentK = 1;
const FS = { root: 10, rim: 12.5, end: 11, inner: 11 };
const fontSize = n => n.data.event && n.parent ? 9.5 : FS[n.labelKind];

// Where a label may sit, in order of preference. Rim and extinction labels have one home;
// branch points try beside, inside and outside the node and keep the first that is clear.
function labelOptions(n) {
  if (n.labelKind === "rim" || n.labelKind === "end") return ["radial"];
  if (!n.parent) return ["right", "left", "below"];
  return n.data.event ? ["right", "left", "in", "out"] : ["in", "out", "right", "left"];
}

// Label geometry for a placement mode at polar position (t, r). dy is in em.
function labelGeom(n, mode, t, r) {
  const g = { x: 0, y: 0, rot: 0, anchor: "middle", dy: 0.35 };
  const deg = t * 180 / Math.PI, a = norm(t);
  if (mode === "radial") {
    const rl = n.labelKind === "rim" ? LEAF_LABEL_R : Math.max(r + n.nr + 4, TIME(n.lastEnd) + 5);
    const flip = a > Math.PI;
    [g.x, g.y] = pt(t, rl);
    g.rot = flip ? deg + 90 : deg - 90;
    g.anchor = flip ? "end" : "start";
  } else if (mode === "in" || mode === "out") {
    const top = a < Math.PI / 2 || a > 1.5 * Math.PI;
    [g.x, g.y] = pt(t, mode === "in" ? r - n.nr - 4 : r + n.nr + 4);
    g.rot = top ? deg : deg - 180;
    g.dy = (mode === "in") === top ? 0.8 : 0;
  } else {
    const [x, y] = pt(t, r);
    if (mode === "below") { g.x = x; g.y = y + n.nr + 11; }
    else {
      g.x = x + (mode === "right" ? 1 : -1) * (n.nr + 4);
      g.y = y;
      g.anchor = mode === "right" ? "start" : "end";
    }
  }
  return g;
}

function placeLabel(n, op) {
  const c = n.cur, el = n.labelText, g = labelGeom(n, n.lmode, c.t, c.r);
  const alpha = n.labelKind === "root" ? op : n.labelKind === "inner" ? op * (1 - c.g) : op * c.g;
  el.setAttribute("transform", `translate(${g.x.toFixed(2)},${g.y.toFixed(2)}) rotate(${g.rot.toFixed(2)}) scale(${labelScale.toFixed(3)})`);
  el.setAttribute("text-anchor", g.anchor);
  el.setAttribute("dy", g.dy + "em");
  el.setAttribute("opacity", alpha.toFixed(3));
}

function drawFrame() {
  for (const n of ALL) {
    if (!n.shown) continue;
    const c = n.cur, op = c.a * n.dimC;
    const [x, y] = pt(c.t, c.r);
    n.nodeEl.setAttribute("transform", `translate(${x.toFixed(2)},${y.toFixed(2)})`);
    n.nodeEl.setAttribute("opacity", op.toFixed(3));
    if (n.limbEl) {
      const pts = limbOf(n);
      const w0 = width(n.mass), w1 = Math.max(0.6, w0 * 0.72);
      n.limbEl.setAttribute("d", n.ribbon ? ribbonPath(pts, w0, w1) : centerLine(pts));
      n.limbEl.setAttribute("opacity", op.toFixed(3));
    }
    if (c.g > 0.002) {
      const end = n.data.extinct ? TIME(n.lastEnd) : LEAF_BASE;
      const [xe, ye] = pt(c.t, Math.max(c.r, end));
      const [x2, y2] = pt(c.t, c.r + Math.max(0, end - c.r) * c.g);
      n.tailEl.setAttribute("d", `M${x.toFixed(2)},${y.toFixed(2)}L${x2.toFixed(2)},${y2.toFixed(2)}`);
      n.tailEl.setAttribute("opacity", op.toFixed(3));
      n.tailEl.style.display = "";
      if (n.gradEl) {
        n.gradEl.setAttribute("x1", x); n.gradEl.setAttribute("y1", y);
        n.gradEl.setAttribute("x2", xe); n.gradEl.setAttribute("y2", ye);
      }
    } else {
      n.tailEl.style.display = "none";
    }
    if (n.leafEl) {
      const unfurl = clamp01((c.g - 0.6) / 0.4);          // the leaf opens as its tail arrives
      if (unfurl > 0.001) {
        const [lx, ly] = pt(c.t, LEAF_BASE);
        const s = n.leafScale * d3.easeBackOut(unfurl);
        n.leafEl.setAttribute("transform", `translate(${lx.toFixed(2)},${ly.toFixed(2)}) rotate(${(c.t * 180 / Math.PI + n.tilt).toFixed(2)}) scale(${s.toFixed(3)})`);
        n.leafEl.setAttribute("opacity", op.toFixed(3));
        n.leafEl.style.display = "";
      } else {
        n.leafEl.style.display = "none";
      }
    }
    placeLabel(n, op);
  }
}

// A watercolour wash under the living limbs, laid down once the tree has come to rest.
function drawWash() {
  for (const n of ALL) {
    if (!n.washEl) continue;
    if (!n.visible) { n.washEl.style.display = "none"; continue; }
    const w0 = width(n.mass);
    n.washEl.style.display = "";
    n.washEl.setAttribute("d", ribbonPath(limbOf(n), w0 * 2.8 + 3, w0 * 1.5 + 2));
  }
}

function setShown(n, on) {
  n.shown = on;
  const v = on ? "" : "none";
  n.nodeEl.style.display = v;
  n.labelWrap.style.display = v;
  if (n.limbEl) n.limbEl.style.display = v;
  if (!on) {
    n.tailEl.style.display = "none";
    if (n.leafEl) n.leafEl.style.display = "none";
  }
}

/* ---------- animation ---------- */
let timer = null, glinted = false;
const isWithin = (n, a) => { for (let m = n; m; m = m.parent) if (m === a) return true; return false; };
const clamp01 = v => v < 0 ? 0 : v > 1 ? 1 : v;

function assignKinds() {
  for (const n of ALL) {
    const terminal = n.visible && !n.children;
    n.labelKind = !n.parent ? "root" : !terminal ? "inner" : n.data.extinct ? "end" : "rim";
    n.labelText.setAttribute("class",
      `lbl k-${n.labelKind}` + (n.data.event ? " event" : "") + (n.data.extinct ? " ext" : "") + (!n.parent ? " root" : ""));
  }
}

function update(source, intro) {
  layout();
  assignKinds();
  for (const n of ALL) {
    n.from = { ...n.cur };
    n.to = { t: n.tt, r: n.rt, a: n.visible ? 1 : 0, g: n.visible && !n.children && n.parent ? 1 : 0 };
    n.delay = intro ? n.depth * INTRO_STAGGER
      : source && n.visible && isWithin(n, source) ? (n.depth - source.depth) * STAGGER : 0;
    if (n.visible) setShown(n, true);
    n.nodeEl.classList.toggle("collapsed", !!(n.collapsed && n.allKids));
    if (n.allKids) n.nodeEl.setAttribute("aria-expanded", n.collapsed ? "false" : "true");
  }
  cull();
  gInscFam.classed("faded", true);
  gWash.classed("resting", false);
  run(intro ? INTRO_DUR : DUR);
}

function run(dur) {
  if (timer) timer.stop();
  const total = dur + d3.max(ALL, n => n.delay);
  const ease = d3.easeCubicInOut;
  const step = elapsed => {
    for (const n of ALL) {
      const f = n.from, to = n.to, c = n.cur;
      const u = dur ? clamp01((elapsed - n.delay) / dur) : 1;
      const e = ease(u);
      c.t = f.t + (to.t - f.t) * e;
      c.r = f.r + (to.r - f.r) * e;
      c.a = f.a + (to.a - f.a) * e;
      const ug = to.g > f.g ? clamp01((u - 0.35) / 0.65) : clamp01(u / 0.5);   // tails grow late, retract early
      c.g = f.g + (to.g - f.g) * ease(ug);
    }
    drawFrame();
    if (elapsed >= total) {
      if (timer) { timer.stop(); timer = null; }
      for (const n of ALL) if (!n.visible) setShown(n, false);
      drawInscription();
      gInscFam.classed("faded", false);
      drawWash();
      gWash.classed("resting", true);
      if (glintAnim && !glinted && !REDUCED) {          // once, after the band has finished fading in
        glinted = true;
        setTimeout(() => glintAnim.beginElement(), Math.max(0, 4000 - performance.now()));
      }
    }
  };
  if (total <= 0) { step(0); return; }
  timer = d3.timer(step);
}

/* ---------- hover dimming ---------- */
let dimTimer = null;
function setDim(focus) {
  const keep = new Set();
  if (focus) {
    focus.ancestors().forEach(m => keep.add(m));
    focus.descendants().forEach(m => keep.add(m));
  }
  for (const n of ALL) n.dimT = !focus || keep.has(n) ? 1 : 0.24;
  gWash.classed("dim", !!focus);
  if (REDUCED) { for (const n of ALL) n.dimC = n.dimT; drawFrame(); return; }
  if (dimTimer) return;
  dimTimer = d3.timer(() => {
    let moving = false;
    for (const n of ALL) {
      const d = n.dimT - n.dimC;
      if (Math.abs(d) > 0.01) { n.dimC += d * 0.22; moving = true; } else n.dimC = n.dimT;
    }
    if (!timer) drawFrame();
    if (!moving) { dimTimer.stop(); dimTimer = null; drawFrame(); }
  });
}

/* =====================================================================
   8. LABELS  (greedy collision culling at target positions)
   ===================================================================== */
function textWidth(n) {
  const k = n.data.event || !n.parent ? 0.74 : n.data.extinct ? 0.46 : 0.5;
  return labelText(n).length * fontSize(n) * k;
}

// Approximate a placed label by a row of circles along its baseline.
function geomCircles(g, fs, len) {
  const phi = g.rot * Math.PI / 180, ux = Math.cos(phi), uy = Math.sin(phi);
  const off = (g.dy - 0.35) * fs;                     // glyph centre along the text's own "down"
  const ox = -uy * off, oy = ux * off;
  const a0 = g.anchor === "start" ? 0 : g.anchor === "end" ? -len : -len / 2;
  const out = [], rad = fs * 0.5, step = fs * 0.8;
  for (let u = a0; u <= a0 + len + 0.01; u += step) out.push([g.x + ux * u + ox, g.y + uy * u + oy, rad]);
  return out;
}

// A coarse spatial hash so each label is tested only against its neighbours.
function makeGrid() {
  const CELL = 24, cells = new Map();
  const key = (i, j) => i * 4096 + j;
  return {
    add(c) {
      const k = key(Math.floor(c[0] / CELL), Math.floor(c[1] / CELL));
      if (!cells.has(k)) cells.set(k, []);
      cells.get(k).push(c);
    },
    hits(x, y, r, skip) {
      const i0 = Math.floor(x / CELL), j0 = Math.floor(y / CELL);
      for (let i = i0 - 1; i <= i0 + 1; i++) for (let j = j0 - 1; j <= j0 + 1; j++) {
        const list = cells.get(key(i, j));
        if (!list) continue;
        for (const c of list) {
          if (skip && c[3] === skip) continue;
          const dx = x - c[0], dy = y - c[1], rr = r + c[2];
          if (dx * dx + dy * dy < rr * rr) return true;
        }
      }
      return false;
    }
  };
}

function labelPriority(n) {
  if (!n.parent) return 1000;
  if (n.data.event) return 500 + n.data.size;
  if (n.labelKind === "rim") return 300 + n.data.size;
  if (n.labelKind === "end") return 200 + n.data.size;
  return 100 + n.data.size;
}

function isCandidate(n) {
  if (!n.visible) return false;
  if (!n.parent || n.data.event || n.labelKind !== "inner") return true;
  return n.data.size >= 5 || currentK >= 1.6;      // minor branch points wait for zoom or hover
}

function cull() {
  const labels = makeGrid(), dots = makeGrid(), accepted = new Set();
  for (const n of ALL) {
    n.lmode = labelOptions(n)[0];
    if (n.visible) dots.add([...pt(n.tt, n.rt), n.nr + 0.5, n]);
  }
  const clear = (circles, n) => circles.every(([x, y, r]) => !labels.hits(x, y, r) && !dots.hits(x, y, r, n));
  const cand = ALL.filter(isCandidate).sort((a, b) => labelPriority(b) - labelPriority(a));
  for (const n of cand) {
    const fs = fontSize(n) * labelScale, len = textWidth(n) * labelScale;
    for (const mode of labelOptions(n)) {
      const circles = geomCircles(labelGeom(n, mode, n.tt, n.rt), fs, len);
      if (clear(circles, n)) {
        n.lmode = mode;
        circles.forEach(c => labels.add(c));
        accepted.add(n);
        break;
      }
    }
  }
  for (const n of ALL) n.labelWrap.classList.toggle("culled", !accepted.has(n));
  // century labels only where they are clear of everything else
  gWood.selectAll("text.ring-lbl").classed("culled", function (y) {
    const [x, yy] = pt(RING_LABEL_ANGLE, TIME(y));
    const fs = 10 * labelScale, len = String(this.textContent).length * fs * 0.8;
    const circles = geomCircles({ x, y: yy, rot: 0, anchor: "middle", dy: 0.35 }, fs, len);
    if (!clear(circles, null)) return true;
    circles.forEach(c => labels.add(c));
    return false;
  });
}

/* ---------- the inscription: family names and a dedication in the gilt band ---------- */
let INSC_R = R + 200;
const measureEl = document.getElementById("measure");
function measure(s) {
  measureEl.textContent = s;
  return measureEl.getComputedTextLength();
}

function drawInscription() {
  gInscBase.selectAll("*").remove();
  gInscFam.selectAll("*").remove();
  defs.selectAll(".insc-path").remove();

  // runs of living lineages, by family, in angular order
  const ends = ALL.filter(n => n.visible && n.parent && !n.children && !n.data.extinct).sort((a, b) => a.tt - b.tt);
  const runs = [];
  for (const n of ends) {
    const half = (n.x1 - n.x0) / 2, last = runs[runs.length - 1];
    if (last && last.family === n.family) { last.t1 = n.tt + half; last.mass += n.data.size; }
    else runs.push({ family: n.family, t0: n.tt - half, t1: n.tt + half, mass: n.data.size });
  }
  // each family is named once, over its weightiest run; the bare roots side carries the dates
  const best = new Map();
  for (const run of runs) {
    if (run.family === "trunk") continue;
    if (!best.has(run.family) || best.get(run.family).mass < run.mass) best.set(run.family, run);
  }
  const items = [...best.values()].sort((a, b) => b.mass - a.mass)
    .map(run => ({ text: FAMILIES[run.family].toUpperCase(), t: (run.t0 + run.t1) / 2, run }));
  items.push({ text: "ANNO DOMINI XLIX · MDCC", t: Math.PI });
  const angDist = (a, b) => { const d = Math.abs(norm(a) - norm(b)); return Math.min(d, TAU - d); };
  const placed = [];
  for (const it of items) {
    it.hw = (measure(it.text) / 2 + (it.run ? 16 : 14)) / INSC_R;    // family names open with a larger initial
    const home = it.t;
    for (const shift of [0, .03, -.03, .06, -.06, .1, -.1, .15, -.15, .2, -.2]) {
      it.t = home + shift;
      if (placed.every(p => angDist(p.t, it.t) > p.hw + it.hw)) { placed.push(it); break; }
    }
  }

  // the circle itself, broken where the lettering sits, with a small dot at each break
  const gaps = placed.map(p => [norm(p.t) - p.hw, norm(p.t) + p.hw]).sort((a, b) => a[0] - b[0]);
  gaps.forEach((gap, i) => {
    const next = gaps[(i + 1) % gaps.length];
    const a = gap[1], b = next[0] + (i + 1 === gaps.length ? TAU : 0);
    if (b - a > 0.004) gInscBase.append("path").attr("class", "insc-circle").attr("d", arcPath(a, b, INSC_R, false));
    for (const t of [gap[0], gap[1]]) gInscBase.append("circle").attr("class", "insc-dot").attr("r", 1.6)
      .attr("cx", pt(t, INSC_R)[0]).attr("cy", pt(t, INSC_R)[1]);
  });

  placed.forEach((it, i) => {
    const target = it.run ? gInscFam : gInscBase;
    const bottom = (a => a > Math.PI / 2 && a < 1.5 * Math.PI)(norm(it.t));
    const rr = bottom ? INSC_R + INSC_FS * 0.35 : INSC_R - INSC_FS * 0.35;
    const id = "insc-path-" + i;
    defs.append("path").attr("class", "insc-path").attr("id", id).attr("d", arcPath(it.t - it.hw - 0.3, it.t + it.hw + 0.3, rr, bottom));
    const tp = target.append("text").attr("class", "insc-lbl")
      .append("textPath").attr("href", "#" + id).attr("startOffset", "50%").attr("text-anchor", "middle");
    if (!it.run) {
      // the dedication in display capitals, alternately red and blue, with gilt points between words
      let k = 0;
      for (const ch of it.text) tp.append("tspan").attr("class", /[A-Z]/.test(ch) ? (k++ % 2 ? "b" : "r") : "g").text(ch);
      return;
    }
    tp.append("tspan").attr("class", "init " + (i % 2 ? "b" : "r")).text(it.text[0]);
    tp.append("tspan").text(it.text.slice(1));
    // the family's colour runs along the rim over its own lineages, either side of its name
    const color = PAL[it.run.family];
    const spans = [[it.run.t0 + 0.004, it.t - it.hw], [it.t + it.hw, it.run.t1 - 0.004]];
    for (const [a, b] of spans) {
      if (b - a > 0.004) target.append("path").attr("class", "insc-arc").attr("stroke", color).attr("d", arcPath(a, b, INSC_R, false));
    }
    for (const t of [it.t - it.hw, it.t + it.hw]) target.append("circle").attr("r", 1.7).attr("fill", color)
      .attr("cx", pt(t, INSC_R)[0]).attr("cy", pt(t, INSC_R)[1]);
  });
}

/* =====================================================================
   9. INTERACTION
   ===================================================================== */
const tip = document.getElementById("tip");
let tipNode = null, pinned = false, hideTimer = null, showTimer = null, lastPointer = "mouse", focused = null;

function tipHTML(n) {
  const d = n.data;
  const pips = Array.from({ length: 10 }, (_, i) => `<i class="${i < d.size ? "on" : ""}"></i>`).join("");
  const rows = [`<div class="tt-row"><span class="tt-k">Influence</span><span class="pips" role="img" aria-label="${d.size} of 10">${pips}</span></div>`];
  if (!d.extinct && !d.event) rows.push(`<div class="tt-row"><span class="tt-k">Family</span>${TRADITIONS[n.trad].label}</div>`);
  if (n.parent) rows.push(`<div class="tt-row"><span class="tt-k">From</span>${n.parent.data.name}</div>`);
  if (d.uncertain) rows.push(`<div class="tt-unc"><b>?</b> Uncertain: its origins, extent or line of descent are debated.</div>`);
  if (n.collapsed && n.allKids) rows.push(`<div class="tt-hidden">+${n.descCount} folded branches, click to unfold</div>`);
  const link = d.reference ? `<a class="tt-link" href="${d.reference}" target="_blank" rel="noopener">Wikipedia ↗</a>` : "<span></span>";
  const btn = n.allKids && (pinned || lastPointer === "touch")
    ? `<button class="tt-btn" type="button">${n.collapsed ? "Expand" : "Collapse"} branch</button>` : "";
  return `<div class="tt-name">${d.name}</div><div class="tt-when">${whenText(d)}</div>` +
    `<p class="tt-doc">${d.doctrine}</p>${rows.join("")}<div class="tt-foot">${link}${btn}</div>`;
}

function showTip(n, cx, cy) {
  clearTimeout(hideTimer);
  tipNode = n;
  tip.innerHTML = tipHTML(n);
  tip.style.borderTopColor = n.data.event ? PAL.gold1 : colorOf(n);
  tip.classList.add("on");
  tip.setAttribute("aria-hidden", "false");
  const btn = tip.querySelector(".tt-btn");
  if (btn) btn.addEventListener("click", e => { e.stopPropagation(); toggle(n); pinned = true; showTip(n, cx, cy); });
  const w = tip.offsetWidth, h = tip.offsetHeight;
  let x = cx + 16, y = cy + 14;
  if (x + w > innerWidth - 10) x = cx - w - 16;
  if (y + h > innerHeight - 10) y = innerHeight - h - 10;
  tip.style.left = Math.max(10, x) + "px";
  tip.style.top = Math.max(10, y) + "px";
}

function hideTip() {
  tip.classList.remove("on");
  tip.setAttribute("aria-hidden", "true");
  tipNode = null;
  pinned = false;
}

function focusNode(n) {
  if (focused && focused !== n) focused.labelWrap.classList.remove("force");
  focused = n;
  n.labelWrap.classList.add("force");
  setDim(n);
}

function unfocus() {
  if (focused) focused.labelWrap.classList.remove("force");
  focused = null;
  setDim(null);
}

function dismiss() { unfocus(); hideTip(); }

function toggle(n) {
  if (!n.allKids) return;
  n.collapsed = !n.collapsed;
  update(n);
}

function centerOf(n) {
  const b = n.nodeEl.querySelector(".dot").getBoundingClientRect();
  return [b.x + b.width / 2, b.y + b.height / 2];
}

function bind(sel) {
  sel.on("pointerdown", e => { lastPointer = e.pointerType; })
    .on("pointerenter", (e, n) => {
      lastPointer = e.pointerType;
      if (e.pointerType === "touch") return;
      clearTimeout(hideTimer);
      clearTimeout(showTimer);
      showTimer = setTimeout(() => {
        if (pinned && tipNode !== n) return;
        focusNode(n);
        showTip(n, e.clientX, e.clientY);
      }, 90);
    })
    .on("pointerleave", e => {
      clearTimeout(showTimer);
      if (e.pointerType === "touch" || pinned) return;
      hideTimer = setTimeout(dismiss, 250);
    })
    .on("click", (e, n) => {
      e.stopPropagation();
      if (lastPointer === "touch") { pinned = true; focusNode(n); showTip(n, e.clientX, e.clientY); return; }
      if (n.allKids) toggle(n);
      pinned = !n.allKids;           // a leaf's card stays open so its link can be reached
      focusNode(n);
      showTip(n, e.clientX, e.clientY);
    });
}

const nodeSel = gNodes.selectAll("g.node");
bind(nodeSel);
bind(gLabels.selectAll("g.lw"));

nodeSel
  .on("focus", (e, n) => {
    if (n.nodeEl.matches(":focus-visible")) { focusNode(n); showTip(n, ...centerOf(n)); }
  })
  .on("blur", () => { if (!pinned) dismiss(); })
  .on("keydown", (e, n) => {
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault();
      toggle(n);
      showTip(n, ...centerOf(n));
    }
  });

tip.addEventListener("pointerenter", () => clearTimeout(hideTimer));
tip.addEventListener("pointerleave", () => { if (!pinned) hideTimer = setTimeout(dismiss, 250); });
svg.on("click", () => { if (tipNode) dismiss(); });
document.addEventListener("keydown", e => {
  if (e.key === "Escape") { dismiss(); if (document.activeElement && document.activeElement.blur) document.activeElement.blur(); }
});
// the reading hint bows out once someone starts exploring
for (const type of ["pointerdown", "wheel", "keydown"]) {
  addEventListener(type, () => document.body.classList.add("engaged"), { once: true, passive: true });
}

const zoom = d3.zoom().scaleExtent([0.7, 6]).clickDistance(4)
  .on("start", e => { if (e.sourceEvent && e.sourceEvent.type !== "mousedown" && !pinned) dismiss(); })
  .on("zoom", e => {
    viewport.attr("transform", e.transform);
    if (e.transform.k !== currentK) {
      currentK = e.transform.k;
      labelScale = Math.pow(currentK, -0.65);     // labels grow more slowly than the tree
      drawFrame();
      gWood.selectAll("text.ring-lbl").attr("transform", ringLabelTransform);
    }
  })
  .on("end", () => cull());
svg.call(zoom).on("dblclick.zoom", null);

document.getElementById("btnExpand").addEventListener("click", () => {
  for (const n of ALL) n.collapsed = false;
  dismiss(); update(null);
});
document.getElementById("btnCollapse").addEventListener("click", () => {
  for (const n of ALL) n.collapsed = !!(n.allKids && n.parent && !n.data.event);
  dismiss(); update(null);
});
document.getElementById("btnReset").addEventListener("click", () => {
  svg.transition().duration(REDUCED ? 0 : 650).call(zoom.transform, d3.zoomIdentity);
});

/* =====================================================================
   10. LEGEND & OUTLINE
   ===================================================================== */
function buildOutline() {
  const item = n => {
    const d = n.data;
    const name = d.reference ? `<a href="${d.reference}" target="_blank" rel="noopener">${d.name}</a>` : d.name;
    const kids = n.allKids ? `<ul>${n.allKids.map(item).join("")}</ul>` : "";
    return `<li>${name} <span class="ol-meta">${whenText(d)}${d.uncertain ? " · uncertain" : ""}</span>${kids}</li>`;
  };
  document.getElementById("outlineBody").innerHTML = `<ul>${item(root)}</ul>`;
}

// Size the gilt band and the square viewBox to the outermost label of the full tree.
function fitViewBox() {
  let ext = R + 40;
  for (const n of ALL) {
    if (n.labelKind !== "rim" && n.labelKind !== "end") continue;
    const g = labelGeom(n, "radial", n.tt, n.rt), len = textWidth(n) * 1.05;
    const sgn = g.anchor === "end" ? -1 : 1, phi = g.rot * Math.PI / 180;
    ext = Math.max(ext, Math.hypot(g.x + sgn * Math.cos(phi) * len, g.y + sgn * Math.sin(phi) * len));
  }
  INSC_R = ext + 6 + INSC_HALF;
  const m = INSC_R + INSC_HALF + 28;     // room inside the folio frame
  svg.attr("viewBox", `${-m} ${-m} ${2 * m} ${2 * m}`);
}

/* =====================================================================
   11. THE FOLIO  (gilt frame, red ruling, marginal vines, folding title)
   ===================================================================== */
const folioArt = d3.select("#folio-art");
const titleEl = document.getElementById("title"), btnTitle = document.getElementById("btnTitle");
const LEAF_KINDS = ["gilt", "lapis", "rubric", "verdigris"];

function folioLeaf(g, x, y, deg, kind, s) {
  const l = g.append("g").attr("transform", `translate(${x.toFixed(1)},${y.toFixed(1)}) rotate(${deg.toFixed(1)}) scale(${s.toFixed(2)})`);
  l.append("path").attr("class", "f-leaf " + kind).attr("d", IVY);
  l.append("path").attr("class", "f-vein").attr("d", IVY_VEIN);
}

// A hairline spiral, the flourish a pen makes at the end of a spray.
function curlPath(x, y, r, dir, up) {
  const pts = d3.range(0, 2.4 * Math.PI, 0.25).map(a => {
    const rr = r * (1 - a / (2.9 * Math.PI));
    return [x + dir * (r - rr * Math.cos(a)), y + up * rr * Math.sin(a)];
  });
  return d3.line().curve(d3.curveBasis)(pts);
}

// An ink stem winding down the margin, with an ivy leaf on a tendril at every turn
// and a gold bezant on a hairline opposite, as in a Gothic bar-and-ivy border.
function vine(x0, y0, y1, seed) {
  if (y1 - y0 < 170) return;
  const rnd = mulberry32(seed), g = folioArt.append("g").attr("class", "f-vine");
  const A = 5, LAM = 104;
  const X = y => x0 + A * Math.sin((y - y0) / LAM * TAU);
  g.append("path").attr("class", "f-stem").attr("d", d3.line().curve(d3.curveBasis)(d3.range(y0, y1 + 0.1, 4).map(y => [X(y), y])));
  g.append("path").attr("class", "f-curl").attr("d", curlPath(X(y0), y0, 6, -1, -1));
  g.append("path").attr("class", "f-curl").attr("d", curlPath(X(y1), y1, 6, 1, 1));
  let k = Math.floor(rnd() * 4);
  for (let y = y0 + LAM / 4, s = 1; y < y1 - 24; y += LAM / 2, s = -s) {
    const x = X(y);
    const ex = x + s * (12 + rnd() * 5), ey = y - (9 + rnd() * 6), cx = x + s * 3, cy = y - 13;
    g.append("path").attr("class", "f-curl").attr("d", `M${x.toFixed(1)},${y.toFixed(1)}Q${cx.toFixed(1)},${cy.toFixed(1)} ${ex.toFixed(1)},${ey.toFixed(1)}`);
    folioLeaf(g, ex, ey, Math.atan2(ey - cy, ex - cx) * 180 / Math.PI + 90, LEAF_KINDS[k++ % 4], 1 + rnd() * 0.25);
    const yb = y + LAM / 6, xb = X(yb), bx = xb - s * (9 + rnd() * 3), by = yb + 5;
    g.append("path").attr("class", "f-curl").attr("d", `M${xb.toFixed(1)},${yb.toFixed(1)}Q${(xb - s * 6).toFixed(1)},${(yb - 3).toFixed(1)} ${bx.toFixed(1)},${by.toFixed(1)}`);
    g.append("circle").attr("class", "f-bezant").attr("r", 2.1).attr("cx", bx.toFixed(1)).attr("cy", by.toFixed(1));
  }
}

// The gilt frame and red ruling are CSS (.frame), so they follow the viewport on their own;
// the vines are drawn here in page pixels, only where the side margins are wide enough to hold them.
function drawFolio() {
  const W = innerWidth, H = innerHeight, q = 15;
  folioArt.selectAll("*").remove();
  if ((W - Math.min(W, H)) / 2 < 170 || H < 420) return;
  const folded = titleEl.classList.contains("folded");
  const tb = (folded ? btnTitle : titleEl).getBoundingClientRect().bottom;
  const cb = document.querySelector(".controls").getBoundingClientRect().bottom;
  vine(q + 38, Math.max(tb + 26, H * 0.14), H - 72, 11);
  vine(W - q - 38, Math.max(cb + 34, H * 0.14), H - 72, 29);
}

function setFolded(on, remember) {
  titleEl.classList.toggle("folded", on);
  btnTitle.setAttribute("aria-expanded", String(!on));
  btnTitle.setAttribute("aria-label", on ? "Show title" : "Hide title");
  btnTitle.title = on ? "Show title" : "Hide title";
  if (remember) { try { localStorage.setItem("arbor-title-folded", on ? "1" : "0"); } catch (e) { /* storage unavailable */ } }
  drawFolio();
}
btnTitle.addEventListener("click", () => setFolded(!titleEl.classList.contains("folded"), true));
document.getElementById("btnFold").addEventListener("click", () => { setFolded(true, true); btnTitle.focus(); });
let resizeFrame = 0;
addEventListener("resize", () => { cancelAnimationFrame(resizeFrame); resizeFrame = requestAnimationFrame(drawFolio); });

/* ---------- start ---------- */
let startFolded = false;
try { startFolded = localStorage.getItem("arbor-title-folded") === "1"; } catch (e) { /* storage unavailable */ }
setFolded(startFolded, false);
paint();
buildOutline();
layout();
assignKinds();
fitViewBox();
drawWood();
for (const n of ALL) n.cur = { t: n.tt, r: 0, a: 0, g: 0 };   // the tree grows out from its root
update(null, true);

matchMedia("(prefers-color-scheme: dark)").addEventListener("change", () => {
  paint(); drawInscription(); drawFrame(); drawWash();
});
})();
