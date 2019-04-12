//
//  TransactionNotificationTableViewCell.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

enum TransactionNotificationType {
    case warning

    var image: UIImage? {
        switch self {
        case .warning:
            return UIImage(named: "warning")
        }
    }

    var color: UIColor {
        switch self {
        case .warning:
            return PHColors.red
        }
    }
}

final class TransactionNotificationTableViewCell: UITableViewCell {

    // MARK: - Private properties

    private lazy var stackView: UIStackView = {

        let stackView = UIStackView()
        stackView.translatesAutoresizingMaskIntoConstraints = false

        stackView.axis = .horizontal
        stackView.spacing = 20
        stackView.alignment = .center
        stackView.distribution = .equalCentering

        return stackView
    }()

    private lazy var notificationImageView: UIImageView = {

        let imageView = UIImageView()
        imageView.translatesAutoresizingMaskIntoConstraints = false

        return imageView
    }()

    private lazy var notificationLabel: UILabel = {

        let label = UILabel()
        label.font = PHFonts.bold(ofSize: 10)
        label.textColor = .white

        return label
    }()

    // MARK: - Initialization

    override init(style: UITableViewCell.CellStyle, reuseIdentifier: String?) {
        super.init(style: style, reuseIdentifier: reuseIdentifier)

        configure()
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Configuration

    private func configure() {

        addSubview(stackView)
        let margin: CGFloat = 10
        stackView.topAnchor.constraint(equalTo: topAnchor, constant: margin).isActive = true
        stackView.leftAnchor.constraint(equalTo: leftAnchor, constant: margin).isActive = true
        stackView.rightAnchor.constraint(equalTo: rightAnchor, constant: -margin).isActive = true
        stackView.bottomAnchor.constraint(equalTo: bottomAnchor, constant: -margin).isActive = true

        stackView.addArrangedSubview(notificationImageView)
        stackView.addArrangedSubview(notificationLabel)
    }

    func configure(withType type: TransactionNotificationType, andText text: String) {

        notificationImageView.image = type.image
        backgroundColor = type.color
        notificationLabel.text = text
    }
}
