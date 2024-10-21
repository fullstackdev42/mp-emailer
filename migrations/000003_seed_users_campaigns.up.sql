-- MariaDB dump 10.19  Distrib 10.11.6-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: db    Database: mariadb
-- ------------------------------------------------------
-- Server version	10.11.9-MariaDB-ubu2204

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES
(1,'foobar','jonesrussell42@gmail.com','$2a$10$7U0oMJZ0qtKcrJPI0otrXOTczXRfHdYD64JZ6oB2QTluNMSF9zmE6','2024-10-08 16:41:42');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `campaigns`
--

LOCK TABLES `campaigns` WRITE;
/*!40000 ALTER TABLE `campaigns` DISABLE KEYS */;
INSERT INTO `campaigns` VALUES
(3,'Unmarked Burials','Urge for increased funding towards the investigation into missing children and unmarked burials associated with Indian Residential Schools across Canada.','<p>{{First Name}} {{Last Name}}</p><p>{{Address 1}}</p><p>{{City}}, {{Province}}, {{Postal Code}}</p><p>{{Email Address}}</p><p>{{Date}}</p><p><br></p><p>{{MP\'s Name}}</p><p>House of Commons</p><p>Ottawa, ON</p><p>K1A 0A6</p><p><br></p><p>Dear {{MP\'s Name}},</p><p><br></p><p>I am writing to you as a concerned constituent to urge for increased funding towards the investigation into missing children and unmarked burials associated with Indian Residential Schools across Canada. The work to investigate these institutions has only just begun and the Canadian Government needs to ensure that proper funding is available for organizations to complete their comprehensive investigations.</p><p><br></p><p>The painful legacy of the Indian Residential School system has had devastating impacts on generations of Indigenous peoples. It is crucial that we prioritize truth, justice, and reconciliation by ensuring that every child who never returned home is acknowledged and properly laid to rest with dignity in a matter determined by their families and communities.</p><p><br></p><p>I commend the Government of Canada for taking initial steps in supporting the investigations and providing resources to Indigenous communities. However, the scope and scale of these investigations demand further financial commitment and resources to ensure they are conducted thoroughly, respectfully, and in collaboration with impacted communities.</p><p><br></p><p>Therefore, I urge you to advocate for increased funding in the upcoming budget and parliamentary sessions to support:</p><p><br></p><ol><li data-list=\"ordered\"><span class=\"ql-ui\" contenteditable=\"false\"></span>Comprehensive Investigations: Ensuring that every site suspected of containing unmarked graves is thoroughly investigated using state-of-the-art methods and technologies.</li><li data-list=\"ordered\"><span class=\"ql-ui\" contenteditable=\"false\"></span>Community Support and Healing: Providing adequate resources to support the mental health and well-being of survivors, intergenerational survivors, and affected Indigenous communities by allowing for local and national gatherings.</li><li data-list=\"ordered\"><span class=\"ql-ui\" contenteditable=\"false\"></span>Education and Awareness: Promoting public education and awareness initiatives about the history and ongoing impacts of the residential school system, fostering a deeper understanding of Indigenous histories and cultures.</li><li data-list=\"ordered\"><span class=\"ql-ui\" contenteditable=\"false\"></span>Collaboration with Indigenous Leadership: Ensuring that Indigenous communities lead the process of investigation, commemoration, and healing, respecting their traditional knowledge and protocols.</li><li data-list=\"ordered\"><span class=\"ql-ui\" contenteditable=\"false\"></span>Commemoration: Allow and support impacted communities to honour the children who never returned home in a matter respectful to their beliefs.</li></ol><p><br></p><p>By prioritizing these efforts, Canada can take significant steps towards reconciliation, healing historical wounds, and rebuilding trust with Indigenous peoples.</p><p><br></p><p>As my elected representative, I ask for your unwavering support in advocating for increased funding and resources for these crucial initiatives. Please stand with Indigenous communities and ensure that the necessary investments are made to uncover the truth and honor the memories of those who suffered.</p><p><br></p><p>Thank you for your attention to this urgent matter. I look forward to your continued leadership and commitment to justice and reconciliation.</p><p><br></p><p>Sincerely,</p><p>{{First Name}} {{Last Name}}</p><p>{{Address 1}}</p><p>To find your local MP visit www.ourcommons.ca/Members/en/search.</p>',1,'0000-00-00 00:00:00','2024-10-12 20:24:35'),
(4,'Climate Action Now','Advocate for stronger climate policies and increased investment in renewable energy to combat climate change.','<p>{{First Name}} {{Last Name}}</p><p>{{Address 1}}</p><p>{{City}}, {{Province}}, {{Postal Code}}</p><p>{{Email Address}}</p><p>{{Date}}</p><p><br></p><p>{{MP\'s Name}}</p><p>House of Commons</p><p>Ottawa, ON</p><p>K1A 0A6</p><p><br></p><p>Dear {{MP\'s Name}},</p><p><br></p><p>I am writing to you as a concerned constituent regarding the urgent need for stronger climate action in Canada. The effects of climate change are becoming increasingly evident, and it is crucial that we take immediate and decisive action to mitigate its impacts and transition to a sustainable, low-carbon economy.</p><p><br></p><p>I urge you to support and advocate for the following measures:</p><p><br></p><ol><li>Increase Canada\'s emissions reduction targets to align with the latest climate science</li><li>Accelerate the transition to renewable energy sources</li><li>Implement a more robust carbon pricing system</li><li>Invest in green infrastructure and sustainable transportation</li><li>Support programs for energy-efficient buildings and homes</li><li>Protect and restore natural carbon sinks, such as forests and wetlands</li></ol><p><br></p><p>As my elected representative, I ask that you prioritize climate action in your work and push for ambitious policies that will ensure a sustainable future for all Canadians.</p><p><br></p><p>Thank you for your attention to this critical issue. I look forward to hearing about your efforts to address climate change.</p><p><br></p><p>Sincerely,</p><p>{{First Name}} {{Last Name}}</p>',1,'2024-10-13 10:00:00','2024-10-13 10:00:00');
/*!40000 ALTER TABLE `campaigns` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-10-12 20:34:17
